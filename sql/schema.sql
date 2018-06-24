-- name: drop-all
DROP TABLE IF EXISTS user_keys, oauth_users, oauth_tokens, keys, users, reset_tokens, groups, group_users, community_users CASCADE;

-- name: create-users
CREATE TABLE users (
  id                 UUID PRIMARY KEY,
  created            integer NOT NULL,
  updated            integer NOT NULL,
  username           text UNIQUE NOT NULL,
  type               integer NOT NULL,
  password_hash      bytea NOT NULL,
  email              text UNIQUE NOT NULL,
  name               text default '',
  description        text default '',
  home_url           text default '',
  color              text NOT NULL default '#999999',
  thumb_url          text NOT NULL default '',
  profile_url        text NOT NULL default '',
  poster_url         text NOT NULL default '',
  email_confirmed    boolean DEFAULT false,
  is_admin           boolean DEFAULT false,
  current_key        text NOT NULL default '',
  access_token       text UNIQUE NOT NULL,
  deleted            boolean DEFAULT false
);

-- name: create-oauth_tokens
CREATE TABLE oauth_tokens (
  user_id            UUID NOT NULL references users(id),
  service            text NOT NULL DEFAULT '',
  access_token       text NOT NULL DEFAULT '',
  token_type         text NOT NULL DEFAULT 'Bearer',
  refresh_token      text NOT NULL DEFAULT '',
  expiry             timestamp,
  PRIMARY KEY        (user_id, service)
);

-- name: create-keys
CREATE TABLE keys (
  type               text NOT NULL,
  sha_256            bytea PRIMARY KEY,
  created            integer NOT NULL,
  last_seen          integer NOT NULL,
  name               text,
  user_id            UUID NOT NULL,
  public_bytes       bytea NOT NULL,
  private_bytes      bytea,
  deleted            boolean DEFAULT false
);

-- name: create-reset_tokens
CREATE TABLE reset_tokens (
  id                 UUID PRIMARY KEY,
  created            integer NOT NULL,
  updated            integer NOT NULL,
  email              text NOT NULL,
  used               boolean DEFAULT false
);

-- name: create-community_users
CREATE TABLE community_users (
  community_id        UUID NOT NULL references users(id),
  user_id             UUID NOT NULL references users(id),
  -- TODO - this hsoult be UUID references users(id), but then it would be forced to accept null values, which is a problem
  invited_by          text NOT NULL default '',
  role                text NOT NULL default 'member',
  joined              timestamp,
  PRIMARY KEY         (community_id, user_id)
);