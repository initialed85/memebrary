import createClient from "openapi-fetch";
import type { components, paths } from "./api/api";

const client = createClient<paths>({ baseUrl: "/" });

export type GeneratedMeme = components["schemas"]["Meme"];
type TimelineMeme = {
  id: string;
  original_name: string;
  mime_type: string;
  size: number;
  description: string;
  description_status: string;
  description_generated: boolean;
  tags: string[];
  created_at: string;
  updated_at: string;
  sort_order: number;
  deleted?: boolean;
};

type TimelinePage = {
  memes: TimelineMeme[];
  next_cursor: string;
  total: number;
};

type PollPage = {
  memes: TimelineMeme[];
  next_cursor?: string;
  total?: number;
  server_time?: string;
  truncated?: boolean;
};

/**
 * djangolang caches every generated GET response in Redis under a hash of its
 * query parameters, and only invalidates through the CDC stream. A poll that
 * reuses one parameter set replays the same cached body until something in the
 * object graph changes, so polls go through the compatibility list route (which
 * reads PostgreSQL directly and answers `Cache-Control: no-store`) and are also
 * fetched with `cache: "no-store"` for good measure.
 */
const POLL_PAGE_SIZE = 50;
// updated_at is stamped at transaction start, so a transaction that began
// before our snapshot can still commit after it. Rewind the watermark slightly
// so the next tick re-reads that window; merges are idempotent.
const POLL_OVERLAP_MS = 30_000;
const PAGE_ROUTE = "/api/custom/memes";

function responseError(error: unknown, response?: Response) {
  if (error && typeof error === "object" && "error" in error) {
    const value = (error as { error?: unknown }).error;
    if (typeof value === "string") return new Error(value);
  }
  return new Error(
    `Generated API request failed (${response?.status || "unknown"})`,
  );
}

async function getTagIDs(tag: string): Promise<string[]> {
  const { data, error, response } = await client.GET("/api/tags", {
    params: { query: { name__ilike: tag, limit: 100 } },
  });
  if (error) throw responseError(error, response);
  return (data?.objects || [])
    .map((item) => item.id)
    .filter((id): id is string => Boolean(id));
}

async function getMemeIDsForTagIDs(ids: string[]): Promise<string[]> {
  if (ids.length === 0) return [];
  const { data, error, response } = await client.GET("/api/meme-tags", {
    params: { query: { tag_id__in: ids.join(","), limit: 1000 } },
  });
  if (error) throw responseError(error, response);
  return [
    ...new Set(
      (data?.objects || [])
        .map((item) => item.meme_id)
        .filter((id): id is string => Boolean(id)),
    ),
  ];
}

function toTimelineMeme(item: GeneratedMeme): TimelineMeme | null {
  const tags = (item.referenced_by_meme_tag_meme_id_objects || [])
    .map((link) => link.tag_id_object?.name)
    .filter((name): name is string => Boolean(name))
    .sort();
  if (!item.id) return null;
  return {
    id: item.id,
    original_name: item.original_name || "meme",
    mime_type: item.mime_type || "application/octet-stream",
    size: item.size || 0,
    description: item.description || "",
    description_status: item.description_status || "none",
    description_generated: item.description_generated === 1,
    tags,
    created_at: item.created_at || "",
    updated_at: item.updated_at || "",
    sort_order: item.sort_order || 0,
  };
}

/** The compat handler talks to PostgreSQL directly and always sends `no-store`. */
async function getJSON<T>(path: string): Promise<T> {
  const response = await fetch(path, { cache: "no-store" });
  let body: T | null = null;
  try {
    body = (await response.json()) as T;
  } catch {
    // The error below still reports the status when the body is not JSON.
  }
  if (!response.ok) {
    const error = body as unknown as { error?: string };
    throw new Error(error?.error || `Request failed (${response.status})`);
  }
  return body as T;
}

/**
 * Read one meme by id. djangolang's primary-key route answers `depth=3` but
 * leaves the referenced_by links unexpanded, so the single-object read goes
 * through the list route with `id__in`, which expands each link's tag.
 */
export async function fetchMeme(id: string): Promise<TimelineMeme> {
  const { data, error, response } = await client.GET("/api/memes", {
    params: {
      query: { limit: 1, depth: 3, id__in: id, sort_order__desc: "" },
    },
    cache: "no-store",
  });
  if (error) throw responseError(error, response);
  const object = (data?.objects || [])[0];
  if (!object) {
    // The row is gone (deleted here or in another tab); applyMemes turns this
    // tombstone into a removal from the timeline.
    return {
      id,
      original_name: "",
      mime_type: "",
      size: 0,
      description: "",
      description_status: "none",
      description_generated: false,
      tags: [],
      created_at: "",
      updated_at: new Date().toISOString(),
      sort_order: 0,
      deleted: true,
    };
  }
  return toTimelineMeme(object) as TimelineMeme;
}

/**
 * One poll tick: every meme whose metadata changed since `after`, oldest change
 * first, including rows deleted in the window (marked `deleted`). `cursor`
 * continues a page that overflowed. The returned `watermark` is where the next
 * tick resumes: the cursor when the page was truncated, otherwise the database
 * clock rewound by the overlap window, so client clock skew cannot skip rows.
 */
export async function pollChangedMemes({
  after,
  cursor,
  tag,
}: {
  after: string;
  cursor?: string;
  tag: string;
}): Promise<{ watermark: string; cursor: string; memes: TimelineMeme[] }> {
  const query = new URLSearchParams({
    limit: String(POLL_PAGE_SIZE),
    updated_after: after,
  });
  if (cursor) query.set("cursor", cursor);
  if (tag) query.set("tag", tag);
  const page = await getJSON<PollPage>(`${PAGE_ROUTE}?${query}`);
  const memes = page.memes || [];
  if (page.truncated && page.next_cursor) {
    return {
      watermark: after,
      cursor: page.next_cursor,
      memes,
    };
  }
  const serverTime = page.server_time || new Date().toISOString();
  return {
    watermark: new Date(Date.parse(serverTime) - POLL_OVERLAP_MS).toISOString(),
    cursor: "",
    memes,
  };
}

export async function listMemes({
  limit,
  offset,
  tag,
}: {
  limit: number;
  offset: number;
  tag: string;
}): Promise<TimelinePage> {
  const tagIDs = tag ? await getTagIDs(tag) : [];
  const matchingMemeIDs = tag ? await getMemeIDsForTagIDs(tagIDs) : [];
  if (tag && matchingMemeIDs.length === 0) {
    return { memes: [], next_cursor: "", total: 0 };
  }

  const { data, error, response } = await client.GET("/api/memes", {
    params: {
      query: {
        limit,
        offset,
        depth: 3,
        id__in:
          matchingMemeIDs.length > 0 ? matchingMemeIDs.join(",") : undefined,
        sort_order__desc: "",
      },
    },
    // Timeline pages should describe the current database, not a browser
    // cached snapshot of a previous visit.
    cache: "no-store",
  });
  if (error) throw responseError(error, response);

  const objects = data?.objects || [];
  const memes = objects
    .map((item) => toTimelineMeme(item))
    .filter((item): item is TimelineMeme => item !== null);

  // The generated API exposes offset pagination. Keep the existing frontend
  // cursor-shaped state so the rest of the UI remains unchanged.
  const nextCursor =
    objects.length >= limit ? String(offset + objects.length) : "";
  return {
    memes,
    next_cursor: nextCursor,
    total: data?.total_count || memes.length,
  };
}
