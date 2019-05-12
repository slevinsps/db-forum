CREATE EXTENSION IF NOT EXISTS CITEXT;

DROP TABLE IF EXISTS Post cascade;
DROP TABLE IF EXISTS Users cascade;
DROP TABLE IF EXISTS Forum cascade;
DROP TABLE IF EXISTS Thread cascade;
DROP TABLE IF EXISTS Vote cascade;
DROP INDEX IF EXISTS idx_users;
DROP INDEX IF EXISTS idx_forum;
DROP INDEX IF EXISTS idx_thread;
DROP INDEX IF EXISTS idx_post;
DROP INDEX IF EXISTS idx_vote;

CREATE TABLE IF NOT EXISTS Users (
  about    text,
  email    varchar(100) NOT NULL,
  fullname varchar(100) NOT NULL,
  nickname CITEXT NOT NULL PRIMARY KEY COLLATE "C"  
);


CREATE INDEX IF NOT EXISTS idx_users ON Users 
(
  nickname
);

CREATE TABLE IF NOT EXISTS Forum (
  id      bigserial     PRIMARY KEY,
  posts   bigint        default 0,
  slug    CITEXT        NOT NULL unique,
  threads bigint           default 0,
  title   varchar(100)  NOT NULL,
  "user" CITEXT         NOT NULL REFERENCES Users(nickname) ON DELETE CASCADE   
);

CREATE INDEX IF NOT EXISTS idx_forum ON Forum 
(
  slug
);

CREATE TABLE IF NOT EXISTS Thread (
  id        bigserial           PRIMARY KEY,
  author    CITEXT              NOT NULL   REFERENCES Users(nickname) ON DELETE CASCADE,
  created   TIMESTAMPTZ,
  forum     CITEXT              NOT NULL REFERENCES Forum(slug) ON DELETE CASCADE,
  message   text                NOT NULL,
  slug      CITEXT              DEFAULT NULL,
  title     varchar(100)        NOT NULL,
  votes     int  DEFAULT 0      NOT NULL 
);

CREATE INDEX IF NOT EXISTS idx_thread ON Thread 
(
  slug
);

CREATE TABLE IF NOT EXISTS Post (
  path      text                        NOT NULL,
  id        bigserial                   PRIMARY KEY,
  author    CITEXT                      NOT NULL REFERENCES Users(nickname) ON DELETE CASCADE,
  created   TIMESTAMPTZ,
  forum     CITEXT                      NOT NULL REFERENCES Forum(slug),
  is_edited BOOLEAN                     DEFAULT FALSE       NOT NULL,
  message   text                        NOT NULL,
  parent    BIGINT DEFAULT 0            NOT NULL,
  thread    bigserial                   NOT NULL  REFERENCES Thread(id) ON DELETE CASCADE,
  branch    BIGINT                      NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_post ON Post 
(
  author
);

CREATE TABLE IF NOT EXISTS Vote (
  nickname  CITEXT   NOT NULL  REFERENCES Users(nickname)  ON DELETE CASCADE,
  threadId  bigserial NOT NULL  REFERENCES Thread(id)  ON DELETE CASCADE,
  voice     int       NOT NULL,
  UNIQUE (nickname, threadId)
);

CREATE INDEX IF NOT EXISTS idx_vote ON Vote 
(
  nickname,
  threadId
);

CREATE OR REPLACE FUNCTION insertPost()
  RETURNS TRIGGER AS $$
DECLARE 
    parent_branch INTEGER;
    parent_path text;
BEGIN
  SELECT path, branch into parent_path, parent_branch
  FROM Post 
  WHERE id = NEW.parent;
  IF NEW.parent != 0 
  THEN	
    NEW.branch = parent_branch;
  ELSE
    NEW.branch = NEW.id;
  END IF;
  IF parent_path is null 
  THEN
    NEW.path = cast(0 as bit(32)) || cast(NEW.id as bit(32));
  ELSE
    NEW.path = parent_path || cast(NEW.id as bit(32));
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER insertPost
BEFORE INSERT ON Post
FOR EACH ROW EXECUTE PROCEDURE insertPost();
