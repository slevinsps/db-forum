CREATE EXTENSION IF NOT EXISTS CITEXT;
CREATE EXTENSION IF NOT EXISTS pg_trgm;
DROP TABLE IF EXISTS Post cascade;
DROP TABLE IF EXISTS Users cascade;
DROP TABLE IF EXISTS Forum cascade;
DROP TABLE IF EXISTS Thread cascade;
DROP TABLE IF EXISTS Vote cascade;
DROP TABLE IF EXISTS UsersForum cascade;
DROP INDEX IF EXISTS idx_users;
DROP INDEX IF EXISTS idx_forum;
DROP INDEX IF EXISTS idx_thread;
DROP INDEX IF EXISTS idx_post;
DROP INDEX IF EXISTS idx_user_email;
DROP INDEX IF EXISTS idx_thread_slug;
DROP INDEX IF EXISTS idx_thread_forum;
DROP INDEX IF EXISTS idx_thread_author;
DROP INDEX IF EXISTS idx_post_author;
DROP INDEX IF EXISTS idx_post_thread;
DROP INDEX IF EXISTS idx_post_forum;
DROP INDEX IF EXISTS idx_user_nickname_email;
DROP INDEX IF EXISTS idx_vote;
DROP INDEX IF EXISTS idx_post_branch;
DROP INDEX IF EXISTS idx_post_branch_thread_and_parent;
DROP INDEX IF EXISTS index_users_on_nickname;
DROP INDEX IF EXISTS idx_post_thread_parent;
DROP INDEX IF EXISTS idx_usersForum;


CREATE TABLE IF NOT EXISTS Users (
  about    text,
  email    CITEXT NOT NULL,
  fullname varchar(100) NOT NULL,
  nickname CITEXT  NOT NULL PRIMARY KEY  COLLATE "C"
);

CREATE INDEX IF NOT EXISTS idx_user_nickname_email ON Users USING btree  (nickname, email);
CREATE INDEX IF NOT EXISTS idx_user_email ON Users USING btree  (email);


CREATE TABLE IF NOT EXISTS Forum (
  id      bigserial     PRIMARY KEY,
  posts   bigint        default 0,
  slug    CITEXT        NOT NULL UNIQUE,
  threads bigint           default 0,
  title   varchar(100)  NOT NULL,
  "user"  CITEXT        NOT NULL REFERENCES Users(nickname) ON DELETE CASCADE   
);

CREATE TABLE IF NOT EXISTS Thread (
  id        bigserial           PRIMARY KEY,
  author    CITEXT       NOT NULL   REFERENCES Users(nickname) ON DELETE CASCADE,
  created   TIMESTAMPTZ,
  forum     CITEXT              NOT NULL REFERENCES Forum(slug) ON DELETE CASCADE,
  message   text                NOT NULL,
  slug      CITEXT              DEFAULT NULL,
  title     varchar(100)        NOT NULL,
  votes     int  DEFAULT 0      NOT NULL 
);

CREATE INDEX IF NOT EXISTS idx_thread_author ON Thread  USING btree (author);
CREATE INDEX IF NOT EXISTS idx_thread_slug ON Thread  USING btree (slug);
CREATE INDEX IF NOT EXISTS idx_thread_forum ON Thread  USING btree (forum );


CREATE TABLE IF NOT EXISTS Post (
  path      text                        NOT NULL,
  id        bigserial                   PRIMARY KEY,
  author    CITEXT                      NOT NULL REFERENCES Users(nickname) ON DELETE CASCADE,
  created   TIMESTAMPTZ,
  forum     CITEXT                      NOT NULL,
  is_edited BOOLEAN                     DEFAULT FALSE       NOT NULL,
  message   text                        NOT NULL,
  parent    BIGINT DEFAULT 0            NOT NULL,
  thread    bigserial                   NOT NULL  REFERENCES Thread(id) ON DELETE CASCADE,
  branch    BIGINT                      NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_post_author ON Post USING btree (author COLLATE "C");
CREATE INDEX IF NOT EXISTS idx_post_thread ON Post USING btree (thread);
CREATE INDEX IF NOT EXISTS idx_post_forum ON Post USING btree (forum);
CREATE INDEX IF NOT EXISTS idx_post_branch ON Post USING btree (branch);
CREATE INDEX IF NOT EXISTS idx_post_thread_parent ON Post USING btree (thread, parent);

CREATE TABLE IF NOT EXISTS Vote (
  nickname  CITEXT   NOT NULL  REFERENCES Users(nickname)  ON DELETE CASCADE,
  threadId  bigserial NOT NULL  REFERENCES Thread(id)  ON DELETE CASCADE,
  voice     int       NOT NULL,
  voicePrevious int NOT NULL DEFAULT 0
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_vote_nickname_threadId ON Vote USING btree (nickname, threadId);


CREATE TABLE IF NOT EXISTS UsersForum (
  userNickname  CITEXT NOT NULL,
  forum  CITEXT NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_usersForum ON UsersForum USING btree (userNickname, forum);
drop index if exists idx_usersForum_forum;
CREATE  INDEX IF NOT EXISTS idx_usersForum_forum ON UsersForum USING btree (forum);

drop index if exists idx_usersForum_name;
CREATE  INDEX IF NOT EXISTS idx_usersForum_name ON UsersForum USING btree (userNickname COLLATE "C");


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




CREATE OR REPLACE FUNCTION insertThread()
  RETURNS TRIGGER AS $$
BEGIN
  UPDATE Forum SET threads = threads + 1 where slug = NEW.forum;
  INSERT INTO UsersForum(userNickname, forum) VALUES (NEW.author, NEW.forum) ON CONFLICT (userNickname, forum) DO NOTHING;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER insertThread
BEFORE INSERT ON Thread
FOR EACH ROW EXECUTE PROCEDURE insertThread();



CREATE OR REPLACE FUNCTION updateVote()
  RETURNS TRIGGER AS $$
BEGIN
  UPDATE Thread SET votes = votes - OLD.voice + NEW.voice WHERE id = OLD.ThreadId;
  NEW.voicePrevious := OLD.voice;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER updateVote
BEFORE UPDATE ON Vote
FOR EACH ROW EXECUTE PROCEDURE updateVote();


CREATE OR REPLACE FUNCTION insertVote()
  RETURNS TRIGGER AS $$
BEGIN
  UPDATE Thread SET votes = votes + NEW.voice WHERE id = NEW.ThreadId;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER insertVote
AFTER INSERT ON Vote
FOR EACH ROW EXECUTE PROCEDURE insertVote();
