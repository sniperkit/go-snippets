-- name: delete-users
delete from users;
-- name: insert-users
INSERT INTO users
  (id,created,updated,
    username,type,password_hash,email,name,
    description,
    home_url,email_confirmed,is_admin,current_key,
    access_token,deleted,color,thumb_url,profile_url,poster_url)
VALUES
  ('3fe7d2cc-a8dc-4da0-ac37-c3061d067ae7', 1464282748, 1464282748,
    'b5', 1, '\x2432612431302454704144306a4b4138476d584c5a306e2e6936324c2e6647636b57437a533759634e6655597238666e744b5057576759366f6b4f75', 'b5_test@test.email', 'Brendan O''Brien', 
    'Nullam interdum, lorem ut porttitor ullamcorper, ligula arcu pretium dui, eu feugiat ipsum diam sit amet arcu. Morbi eleifend id orci vitae pulvinar',
    'https://test_website.org', false, false, 'b5_644b51b9567d0d999e40f697d7406a26030cde95a83775d285ff1f57a73b3ebc', 
    'b5', false, '#AAAAAA', 'https://avatars0.githubusercontent.com/u/1154390?v=4&s=100', 'https://avatars0.githubusercontent.com/u/1154390?v=4&s=400', ''),
  ('54b80e91-cae0-423d-b5d8-c9acbb5e2536', 1463687282, 1463687793,
    'ebarry', 1, '\x2432612431302454704144306a4b4138476d584c5a306e2e6936324c2e6647636b57437a533759634e6655597238666e744b5057576759366f6b4f75', 'ebarry_test@test.email', 'ebarry',
    'Nullam interdum, lorem ut porttitor ullamcorper, ligula arcu pretium dui, eu feugiat ipsum diam sit amet arcu. Morbi eleifend id orci vitae pulvinar',
    'https://test_website.org', false, false, 'ebarry_644b51b9567d0d999e40f697d7406a26030cde95a83775d285ff1f57a73b3ebc', 
    'ebarry', false, '#AAAAAA', 'https://avatars6.githubusercontent.com/u/161439?v=4&s=100', 'https://avatars6.githubusercontent.com/u/161439?v=4&s=400', ''),
  ('0232fb99-f965-4fe5-bec9-ad099760ab29', 1464282748, 1464282748, 
    'blackglade', 1, '\x2432612431302454704144306a4b4138476d584c5a306e2e6936324c2e6647636b57437a533759634e6655597238666e744b5057576759366f6b4f75', 'blackglade_test@test.email', 'Harsh Baid',
    'Nullam interdum, lorem ut porttitor ullamcorper, ligula arcu pretium dui, eu feugiat ipsum diam sit amet arcu. Morbi eleifend id orci vitae pulvinar', 
    'https://test_website.org', false, false, 'blackglade_644b51b9567d0d999e40f697d7406a26030cde95a83775d285ff1f57a73b3ebc', 
    'blackglade', false, '#AAAAAA', 'https://avatars0.githubusercontent.com/u/24626991?v=3&s=100', 'https://avatars0.githubusercontent.com/u/24626991?v=3&s=400', ''),
  ('1b674f47-d0f4-4b3c-b25d-c49521b5599a', 1464282748, 1464282748, 
    'murphyofglad', 1, '\x2432612431302454704144306a4b4138476d584c5a306e2e6936324c2e6647636b57437a533759634e6655597238666e744b5057576759366f6b4f75', 'murphyofglad_test@test.email', 'Michelle Murphy',
    'Lorem ipsum dolor sit amet, consectetur adipiscing elit. In fringilla vulputate justo eu tempus. Sed eleifend dui massa, in accumsan leo convallis id. Proin vitae leo lacus. Nunc rhoncus augue quis iaculis faucibus. Cras sapien purus, pulvinar non sodales vitae, congue ut dui. Sed egestas velit a est placerat mollis.', 
    'https://test_website.org', false, false, 'murphyofglad_644b51b9567d0d999e40f697d7406a26030cde95a83775d285ff1f57a73b3ebc', 
    'murphyofglad', false, '#AAAAAA', 'https://avatars0.githubusercontent.com/u/24626991?v=3&s=100', 'https://avatars0.githubusercontent.com/u/24626991?v=3&s=400', ''),
  ('a5af4a22-746c-4c10-a74f-99a2867f96fb', 1492536378, 1499898140, 
    'jeffliu', 1, '\x2432612431302454704144306a4b4138476d584c5a306e2e6936324c2e6647636b57437a533759634e6655597238666e744b5057576759366f6b4f75', 'jeffliu_test@test.email', 'Jeffrey Liu', 
    'Lorem ipsum dolor sit amet, consectetur adipiscing elit. In fringilla vulputate justo eu tempus. Sed eleifend dui massa, in accumsan leo convallis id. Proin vitae leo lacus. Nunc rhoncus augue quis iaculis faucibus. Cras sapien purus, pulvinar non sodales vitae, congue ut dui. Sed egestas velit a est placerat mollis.', 
    'https://test_website.org', false, false, 'jeffliu_644b51b9567d0d999e40f697d7406a26030cde95a83775d285ff1f57a73b3ebc', 
    'jeffliu', false, '#AAAAAA', 'https://avatars3.githubusercontent.com/u/1486277?v=3&s=150', 'https://avatars3.githubusercontent.com/u/1486277?v=3&s=400', ''),
  ('57013bf0-2366-11e6-b67b-9e71128cae77', 1464282748, 1499908984, 
    'EDGI', 2, '\x2432612431302454704144306a4b4138476d584c5a306e2e6936324c2e6647636b57437a533759634e6655597238666e744b5057576759366f6b4f75', 'edgi_test@test.email', 'Environmental Data & Governance Initiative', 
    'The Environmental Data & Governance Initiative (EDGI) is an international network of academics and non-profits addressing potential threats to federal environmental and energy policy, and to the scientific research infrastructure built to investigate, inform, and enforce. Dismantling this infrastructure -- which ranges from databases to satellites to models for climate, air, and water -- could imperil the publics right to know, the United States standing as a scientific leader, corporate accountability, and environmental protection.', 
    'https://envirodatagov.org', false, false, 'EDGI_644b51b9567d0d999e40f697d7406a26030cde95a83775d285ff1f57a73b3ebc', 
    'EDGI', false, '#AAAAAA', 'https://avatars0.githubusercontent.com/u/24626991?v=3&s=200', 'https://avatars0.githubusercontent.com/u/24626991?v=3&s=400', '');


-- name: delete-reset_tokens
delete from reset_tokens;
-- name: insert-reset_tokens
INSERT INTO reset_tokens VALUES
  ('69eb9cbd-7085-4624-a841-59d0f02eaa7b', 1464282748, 1464282748, 'test_user_brendan@qri.io', false);

-- name: delete-keys
delete from keys;
-- name: insert-keys
INSERT INTO keys 
  (user_id, name, type, sha_256, created, last_seen, public_bytes, private_bytes, deleted)
VALUES
  -- TODO - currently these are just placeholders, keys feature needs work
  -- ('a5af4a22-746c-4c10-a74f-99a2867f96fb', 'default key', 'ssh', '\xac4c6382d0f7c494570affb1767e4da8d55d39a011a25f2360832a1e1260b21b', 1492536378, 0, '\x000000077373682d727361000000030100010000008100ae05019d0e3c65aa8b868d1ca62ed4ef56c08e230967d30d025724d596c8eae76cd5eb31e9adecd902a2311714f116e88d5c3da6d7263a','\x',false),
  -- ('286e9d09-1ba9-42fd-8de8-466cc131ea58', 'default key', 'ssh', '\x14f216ed95448ab26627e7bb37142a36c50087dee727cf2dc6088a52641edb8c', 1499991357, 0, '\x000000077373682d727361000000030100010000008100dc6a91e36fa1cc047016d18779146c816d3cfb3425aa70574c1d55154bed850404a1c13f4226c24c66681376dfebfd93a43901c9f1a769', '\x', false);
  ('3fe7d2cc-a8dc-4da0-ac37-c3061d067ae7', 'default key', 'ssh', '\x114c6382d0f7c494570affb1767e4da8d55d39a011a25f2360832a1e1260b21b', 1492536378, 0, '\x000000077373682d727361000000030100010000008100ae05019d0e3c65aa8b868d1ca62ed4ef56c08e230967d30d025724d596c8eae76cd5eb31e9adecd902a2311714f116e88d5c3da6d7263a','\x',false),
  ('54b80e91-cae0-423d-b5d8-c9acbb5e2536', 'default key', 'ssh', '\x22f216ed95448ab26627e7bb37142a36c50087dee727cf2dc6088a52641edb8c', 1492536378, 0, '\x000000077373682d727361000000030100010000008100ae05019d0e3c65aa8b868d1ca62ed4ef56c08e230967d30d025724d596c8eae76cd5eb31e9adecd902a2311714f116e88d5c3da6d7263a','\x',false),
  ('0232fb99-f965-4fe5-bec9-ad099760ab29', 'default key', 'ssh', '\x334c6382d0f7c494570affb1767e4da8d55d39a011a25f2360832a1e1260b21b', 1492536378, 0, '\x000000077373682d727361000000030100010000008100ae05019d0e3c65aa8b868d1ca62ed4ef56c08e230967d30d025724d596c8eae76cd5eb31e9adecd902a2311714f116e88d5c3da6d7263a','\x',false),
  ('1b674f47-d0f4-4b3c-b25d-c49521b5599a', 'default key', 'ssh', '\x444c6382d0f7c494570affb1767e4da8d55d39a011a25f2360832a1e1260b21b', 1492536378, 0, '\x000000077373682d727361000000030100010000008100ae05019d0e3c65aa8b868d1ca62ed4ef56c08e230967d30d025724d596c8eae76cd5eb31e9adecd902a2311714f116e88d5c3da6d7263a','\x',false),
  ('a5af4a22-746c-4c10-a74f-99a2867f96fb', 'default key', 'ssh', '\x554c6382d0f7c494570affb1767e4da8d55d39a011a25f2360832a1e1260b21b', 1492536378, 0, '\x000000077373682d727361000000030100010000008100ae05019d0e3c65aa8b868d1ca62ed4ef56c08e230967d30d025724d596c8eae76cd5eb31e9adecd902a2311714f116e88d5c3da6d7263a','\x',false),
  ('57013bf0-2366-11e6-b67b-9e71128cae77', 'default key', 'ssh', '\x664c6382d0f7c494570affb1767e4da8d55d39a011a25f2360832a1e1260b21b', 1492536378, 0, '\x000000077373682d727361000000030100010000008100ae05019d0e3c65aa8b868d1ca62ed4ef56c08e230967d30d025724d596c8eae76cd5eb31e9adecd902a2311714f116e88d5c3da6d7263a','\x',false);


-- name: delete-oauth_tokens
DELETE FROM oauth_tokens;
-- name: insert-oauth_tokens
-- INSERT INTO oauth_tokens VALUES ();

-- name: delete-community_users
DELETE FROM community_users;
-- name: insert-community_users
INSERT INTO community_users
  (community_id, user_id, invited_by, role, joined)
VALUES 
  ('57013bf0-2366-11e6-b67b-9e71128cae77', 'a5af4a22-746c-4c10-a74f-99a2867f96fb', '', 'admin', '2017-03-23 00:00:01'),
  ('57013bf0-2366-11e6-b67b-9e71128cae77', '0232fb99-f965-4fe5-bec9-ad099760ab29', 'a5af4a22-746c-4c10-a74f-99a2867f96fb', 'member', '2017-03-23 00:00:01'),
  ('57013bf0-2366-11e6-b67b-9e71128cae77', '1b674f47-d0f4-4b3c-b25d-c49521b5599a', 'a5af4a22-746c-4c10-a74f-99a2867f96fb', 'member', '2017-03-23 00:00:01');

