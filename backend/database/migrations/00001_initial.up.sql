CREATE SCHEMA IF NOT EXISTS public;

SET LOCAL search_path = public;

--
-- meme
--

CREATE TABLE
    public.meme (
        id uuid PRIMARY KEY NOT NULL UNIQUE DEFAULT gen_random_uuid (),
        created_at timestamptz NOT NULL DEFAULT now(),
        updated_at timestamptz NOT NULL DEFAULT now(),
        deleted_at timestamptz NULL DEFAULT NULL,
        filename TEXT NOT NULL UNIQUE,
        original_name TEXT NOT NULL,
        mime_type TEXT NOT NULL,
        size INTEGER NOT NULL,
        description TEXT NOT NULL DEFAULT '',
        description_status TEXT NOT NULL DEFAULT 'none',
        description_generated INTEGER NOT NULL DEFAULT 0,
        sort_order INTEGER NOT NULL DEFAULT 0,
        metadata_version INTEGER NOT NULL DEFAULT 0,
        ai_worker_claimed_until timestamptz NOT NULL DEFAULT to_timestamp(0)
    );

ALTER TABLE public.meme OWNER TO postgres;

CREATE
OR REPLACE FUNCTION create_meme () RETURNS TRIGGER AS $$
BEGIN
  NEW.created_at = now();
  NEW.updated_at = now();
  NEW.deleted_at = null;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER create_meme BEFORE INSERT ON public.meme FOR EACH ROW
EXECUTE PROCEDURE create_meme ();

CREATE
OR REPLACE FUNCTION update_meme () RETURNS TRIGGER AS $$
BEGIN
  NEW.created_at = OLD.created_at;
  NEW.updated_at = now();
  IF OLD.deleted_at IS NOT null AND NEW.deleted_at IS NOT null THEN
    NEW.deleted_at = OLD.deleted_at;
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_meme BEFORE
UPDATE ON public.meme FOR EACH ROW
EXECUTE PROCEDURE public.update_meme ();

CREATE RULE delete_meme AS ON DELETE TO public.meme
DO INSTEAD (
    UPDATE public.meme
    SET
        created_at = old.created_at,
        updated_at = now(),
        deleted_at = now()
    WHERE
        id = old.id
        AND deleted_at IS null
);

--
-- tag
--

CREATE TABLE
    public.tag (
        id uuid PRIMARY KEY NOT NULL UNIQUE DEFAULT gen_random_uuid (),
        created_at timestamptz NOT NULL DEFAULT now(),
        updated_at timestamptz NOT NULL DEFAULT now(),
        deleted_at timestamptz NULL DEFAULT NULL,
        name TEXT NOT NULL UNIQUE
    );

ALTER TABLE public.tag OWNER TO postgres;

CREATE
OR REPLACE FUNCTION create_tag () RETURNS TRIGGER AS $$
BEGIN
  NEW.created_at = now();
  NEW.updated_at = now();
  NEW.deleted_at = null;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER create_tag BEFORE INSERT ON public.tag FOR EACH ROW
EXECUTE PROCEDURE create_tag ();

CREATE
OR REPLACE FUNCTION update_tag () RETURNS TRIGGER AS $$
BEGIN
  NEW.created_at = OLD.created_at;
  NEW.updated_at = now();
  IF OLD.deleted_at IS NOT null AND NEW.deleted_at IS NOT null THEN
    NEW.deleted_at = OLD.deleted_at;
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_tag BEFORE
UPDATE ON public.tag FOR EACH ROW
EXECUTE PROCEDURE public.update_tag ();

CREATE RULE delete_tag AS ON DELETE TO public.tag
DO INSTEAD (
    UPDATE public.tag
    SET
        created_at = old.created_at,
        updated_at = now(),
        deleted_at = now()
    WHERE
        id = old.id
        AND deleted_at IS null
);

--
-- tag -> meme (many-to-many)
--

CREATE TABLE
    public.meme_tag (
        id uuid PRIMARY KEY NOT NULL UNIQUE DEFAULT gen_random_uuid (),
        created_at timestamptz NOT NULL DEFAULT now(),
        updated_at timestamptz NOT NULL DEFAULT now(),
        deleted_at timestamptz NULL DEFAULT NULL
    );

ALTER TABLE public.meme_tag OWNER TO postgres;

CREATE
OR REPLACE FUNCTION create_meme_tag () RETURNS TRIGGER AS $$
BEGIN
  NEW.created_at = now();
  NEW.updated_at = now();
  NEW.deleted_at = null;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER create_meme_tag BEFORE INSERT ON public.meme_tag FOR EACH ROW
EXECUTE PROCEDURE create_meme_tag ();

CREATE
OR REPLACE FUNCTION update_meme_tag () RETURNS TRIGGER AS $$
BEGIN
  NEW.created_at = OLD.created_at;
  NEW.updated_at = now();
  IF OLD.deleted_at IS NOT null AND NEW.deleted_at IS NOT null THEN
    NEW.deleted_at = OLD.deleted_at;
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_meme_tag BEFORE
UPDATE ON public.meme_tag FOR EACH ROW
EXECUTE PROCEDURE public.update_meme_tag ();

CREATE RULE delete_meme_tag AS ON DELETE TO public.meme_tag
DO INSTEAD (
    UPDATE public.meme_tag
    SET
        created_at = old.created_at,
        updated_at = now(),
        deleted_at = now()
    WHERE
        id = old.id
        AND deleted_at IS null
);


ALTER TABLE public.meme_tag
ADD COLUMN meme_id uuid NOT NULL REFERENCES public.meme (id);


ALTER TABLE public.meme_tag
ADD COLUMN tag_id uuid NOT NULL REFERENCES public.tag (id);


