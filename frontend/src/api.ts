import createClient from "openapi-fetch";
import type { components, paths } from "./api/api";

const client = createClient<paths>({ baseUrl: "/" });

export type GeneratedMeme = components["schemas"]["Meme"];
export type GeneratedMemeTag = components["schemas"]["MemeTag"];

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
};

type TimelinePage = {
  memes: TimelineMeme[];
  next_cursor: string;
  total: number;
};

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

async function getTagsByMemeIDs(ids: string[]): Promise<Map<string, string[]>> {
  const tagsByMeme = new Map<string, string[]>();
  if (ids.length === 0) return tagsByMeme;

  const { data, error, response } = await client.GET("/api/meme-tags", {
    params: {
      query: {
        meme_id__in: ids.join(","),
        // meme-tag is the direct endpoint for this relationship. Asking the
        // meme endpoint for depth=3 also loads meme-tag -> meme, duplicating
        // the parent graph; tag__load gives us only the useful second hop.
        tag__load: "",
        limit: Math.max(100, ids.length * 16),
      },
    },
  });
  if (error) throw responseError(error, response);

  for (const item of data?.objects || []) {
    const memeID = item.meme_id;
    const tag = item.tag_id_object?.name;
    if (!memeID || !tag) continue;
    const existing = tagsByMeme.get(memeID) || [];
    existing.push(tag);
    tagsByMeme.set(memeID, existing);
  }
  for (const tags of tagsByMeme.values()) tags.sort();
  return tagsByMeme;
}

function toTimelineMeme(
  item: GeneratedMeme,
  tags: string[],
): TimelineMeme | null {
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
        id__in:
          matchingMemeIDs.length > 0 ? matchingMemeIDs.join(",") : undefined,
        sort_order__desc: "",
      },
    },
  });
  if (error) throw responseError(error, response);

  const objects = data?.objects || [];
  const ids = objects
    .map((item) => item.id)
    .filter((id): id is string => Boolean(id));
  const tagsByMeme = await getTagsByMemeIDs(ids);
  const memes = objects
    .map((item) => toTimelineMeme(item, tagsByMeme.get(item.id || "") || []))
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
