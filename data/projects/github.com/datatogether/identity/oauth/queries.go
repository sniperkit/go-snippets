package oauth

const qUserOauthTokenUser = `
SELECT
  id, created, updated, username, type, name, description, home_url, email, current_key, email_confirmed, is_admin
FROM users
WHERE
  id = (SELECT user_id FROM oauth_tokens WHERE service = $1 AND token = $2);`

// const qUserOauthTokens = `
// SELECT
// FROM oauth_tokens;`

const qUserOauthTokenInsert = `
INSERT INTO oauth_tokens 
  (user_id, service, access_token, token_type, refresh_token, expiry) 
VALUES 
  ($1, $2, $3, $4, $5, $6);`

const qUserOauthTokenUpdate = `
UPDATE oauth_tokens SET
  user_id = $1, service = $2, access_token = $3, token_type = $4, refresh_token = $5, expiry = $6
WHERE
  user_id = $1 AND
  service = $2;`

const qUserOauthTokenDelete = `DELETE FROM oauth_tokens WHERE user_id = $1 AND service = $2;`

const qUserOauthTokenByAccessToken = `
SELECT
  user_id, service, access_token, token_type, refresh_token, expiry
FROM oauth_tokens
WHERE 
  access_token = $1;`

const qUserOauthTokenByUserAndService = `
SELECT
  user_id, service, access_token, token_type, refresh_token, expiry
FROM oauth_tokens
WHERE 
  user_id = $1 AND
  service = $2;`

const qUserOauthServiceToken = `
SELECT
  user_id, service, access_token, token_type, refresh_token, expiry
FROM oauth_tokens
WHERE
  service = $1 AND
  user_id = $2`

const qUserOauthTokensForUser = `
SELECT
  user_id, service, access_token, token_type, refresh_token, expiry
FROM oauth_tokens
WHERE
  user_id = $1;`
