CREATE EXTENSION IF NOT EXISTS CITEXT;

DROP TABLE IF EXISTS Post;
DROP TABLE IF EXISTS Thread;
DROP TABLE IF EXISTS Forum;
DROP TABLE IF EXISTS Vote;
DROP TABLE IF EXISTS Users;

CREATE TABLE IF NOT EXISTS Users (
  about    text,
  email    varchar(100) NOT NULL,
  fullname varchar(100) NOT NULL,
  nickname CITEXT NOT NULL PRIMARY KEY COLLATE "C"  
);
  
CREATE TABLE IF NOT EXISTS Forum (
  id      bigserial PRIMARY KEY,
  posts   bigint       default 0,
  slug    varchar(100) NOT NULL unique,
  threads int          default 0,
  title   varchar(100)  NOT NULL,
  "user" CITEXT NOT NULL REFERENCES Users(nickname)    
);

CREATE TABLE IF NOT EXISTS Thread (
  id        bigserial           PRIMARY KEY,
  author    CITEXT              NOT NULL   REFERENCES Users(nickname),
  created   TIMESTAMPTZ,
  forum     varchar(100)        NOT NULL,
  message   text                NOT NULL,
  slug      varchar(100)        DEFAULT NULL,
  title     varchar(100)        NOT NULL,
  votes     int  DEFAULT 0      NOT NULL 
);

CREATE TABLE IF NOT EXISTS Post (
  path      text                        NOT NULL,
  id        bigserial                   PRIMARY KEY,
  author    CITEXT                      NOT NULL REFERENCES Users(nickname),
  created   TIMESTAMPTZ,
  forum     varchar(100)                NOT NULL,
  is_edited BOOLEAN DEFAULT FALSE       NOT NULL,
  message   text                        NOT NULL,
  parent    BIGINT DEFAULT 0            NOT NULL,
  thread    bigserial                   NOT NULL,
  branch    BIGINT                      NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS Vote (
    nickname  CITEXT   NOT NULL REFERENCES Users(nickname)  unique,
    voice     int       NOT NULL
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
