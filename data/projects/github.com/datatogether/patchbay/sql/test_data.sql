
-- name: insert-primers
insert into primers 
  (id,created,updated,short_title,title,description,parent_id,deleted)
values
  ('5b1031f4-38a8-40b3-be91-c324bf686a87','2017-01-01 00:00:01','2017-01-01 00:00:01', 'EPA', 'Environmental Protection Agency', 'The mission of the Environmental Protection Agency is to protect human health and the environment through the development and enforcement of regulations. The EPA is responsible for administering a number of laws that span various sectors, such as agriculture, transportation, utilities, construction, and oil and gas. In the budget for FY 2017, the agency lays out goals to better support communities and address climate change following the President’s Climate Action Plan. Additionally, the agency aims to improve community water infrastructure, chemical plant safety, and collaborative partnerships among federal, state, and tribal levels.', '', false),
  ('d99891f3-cfd9-4410-aaa4-6e90d792a20a','2017-01-01 00:00:01','2017-01-01 00:00:01', 'Sub-EPA', 'Sub Environmental Protection Agency', 'Sub-Primer for EPA', '5b1031f4-38a8-40b3-be91-c324bf686a87', false),
  ('d9deff9d-15e8-43f1-9d00-51160c0bffbe','2017-01-01 00:00:01','2017-01-01 00:00:01', 'US Census', 'United States Census Bureau', 'US Census Bureau', '',false),
  ('644184d8-042b-4d0d-9bb8-81f142b99dc8','2017-03-23 00:00:01','2017-03-23 00:00:01', 'DOE', 'Department of Energy', 'The mission of the Department of Energy is to ensure America’s security and prosperity by addressing its energy, environmental and nuclear challenges through transformative science and technology solutions. The DOE is broadly responsible for nuclear security, science and energy research, and other regulatory activities. In the budget for FY 2017, the department’s strategic plan is outlined as advancing foundational science, innovating energy technologies, and informing data‐driven policies that enhance U.S. economic growth and job creation, energy security, and environmental quality, with emphasis on implementation of the President’s Climate Action Plan to mitigate the risks of and enhance resilience against climate change.', '', false),
  ('cd5ead69-26a5-4cec-b8e2-90da286a4613','2017-03-23 00:00:01','2017-03-23 00:00:01', 'NASA', 'National Aeronautics and Space Administration', 'NASA’s vision is to “reach for new heights and reveal the unknown for the benefit of humankind”. It is split into 4 mission directorates, which focus on aeronautics, human exploration and operations, science, and space technology, and conducts research at 18 centers and facilities. In 2014, NASA published a strategic plan outlining goals for the next 5-10 years. The goals of the agency are to explore space through both unmanned- and manned-space flight, advance understanding of Earth and its system, and develop technologies to improve life on Earth.', '', false),
  ('0f28f814-ec4c-4a75-b26e-4de435deb218','2017-03-23 00:00:01','2017-03-23 00:00:01', 'NOAA', 'National Oceanic and Atmospheric Administration', 'The mission of the National Oceanic and Atmospheric Administration (NOAA) is focused on science, service, and stewardship, and its three main goals are 1) to understand and predict changes in climate, weather, oceans and coasts; 2) to share that knowledge and information with others; and 3) to conserve and manage coastal and marine ecosystems and resources. NOAA works towards a vision of the future that includes healthy ecosystems, communities, and economies that are resilient to changes. NOAA’s charter is outlined in the Magnus-Stevenson Act Report to Congress.', '', false);
-- name: delete-primers
delete from primers;

--name: insert-sources
insert into sources
  (id,created,updated,title,description,url,primer_id,crawl,stale_duration,last_alert_sent,stats,meta)
values
  ('326fcfa0-d3e6-4b2d-8f95-e77220e16109', '2017-01-01 00:00:01', '2017-01-01 00:00:01', 'epa.gov', 'entire epa site', 'www.epa.gov', '5b1031f4-38a8-40b3-be91-c324bf686a87',true,43200000,null,null,null),
  ('590e001b-7060-4e54-bc81-c20c305a8155', '2017-04-25 00:00:01', '2017-04-25 00:00:01', 'Hazardous Air Pollutants', 'Office of Air and Radiation', 'www.epa.gov/haps', '5b1031f4-38a8-40b3-be91-c324bf686a87',true,43200000,null,null,null),
  ('440d9779-406c-4015-8f2d-404b04ead3a2', '2017-01-01 00:00:01', '2017-01-01 00:00:01', 'census.gov', 'entire census site', 'www.census.gov', 'd9deff9d-15e8-43f1-9d00-51160c0bffbe',true,43200000,null,null,null),
  ('29855324-f444-4c7f-a7a9-936ee4da538a', '2017-01-01 00:00:01', '2017-01-01 00:00:01', 'data.census.gov', 'non-existent census subdomain', 'data.census.gov', 'd9deff9d-15e8-43f1-9d00-51160c0bffbe',false,43200000,null,null,null);
--name: delete-sources
delete from sources;

-- name: insert-urls
insert into urls
  (url,created,updated,last_head,last_get,status,content_type,content_sniff,content_length,file_name,title,id,headers_took,download_took,headers,meta,hash)
values
  ('http://www.epa.gov', '2017-01-01 00:00:01', '2017-01-01 00:00:01', '2017-01-01 00:00:01', null, 200, 'text/html; charset=utf-8', 'text/html;', -1, '', 'United States Environmental Protection Agency, US EPA', 'cee7bbd4-2bf9-4b83-b2c8-be6aeb70e771',0,0, '["X-Content-Type-Options","nosniff","Expires","Fri, 24 Feb 2017 21:53:45 GMT","Date","Fri, 24 Feb 2017 21:53:45 GMT","Etag","W/\"7f53-549471782bb42\"","X-Ua-Compatible","IE=Edge,chrome=1","X-Cached-By","Boost","Content-Type","text/html; charset=utf-8","Vary","Accept-Encoding","Accept-Ranges","bytes","Cache-Control","no-cache, no-store, must-revalidate, post-check=0, pre-check=0","Server","Apache","Connection","keep-alive","Strict-Transport-Security","max-age=31536000; preload;"]', null, '1220459219b10032cc86dcdbc0f83aea15a9d3e1119e7b5170beaee233008ea2c2de'),
  ('https://www.census.gov/nometa.pdf','2017-03-15 17:36:40', '2017-03-21 22:25:21','2017-03-21 22:25:20.88', '2017-03-21 22:25:20.88', 200,'text/html','application/pdf; charset=utf-8' ,164010, 'nometa.pdf','North American Industry Classification System (NAICS) Main Page � U.S. Census Bureau','4c5fc7b8-1397-4d34-980b-1d01247f9ee4',0,0 ,'["Date","Tue, 21 Mar 2017 22:25:20 GMT","Accept-Ranges","bytes","Content-Type","text/html","Strict-Transport-Security","max-age=31536000","Vary","Accept-Encoding"]',null,'1220af06510193276b5fd9ad2fc55dcc004ada557d9259ca3505478bfef0b12ed988'),
  ('https://www.census.gov/topics/economy/classification-codes.html','2017-03-15 17:36:40', '2017-03-21 22:25:21','2017-03-21 22:25:20.88','2017-03-21 22:25:20.88384',200,'text/html','text/plain; charset=utf-8' ,164010, '','North American Industry Classification System (NAICS) Main Page � U.S. Census Bureau','4c5fc7b8-1397-4d34-980b-1d01247f9ee4',0,0 ,'["Date","Tue, 21 Mar 2017 22:25:20 GMT","Accept-Ranges","bytes","Content-Type","text/html","Strict-Transport-Security","max-age=31536000","Vary","Accept-Encoding"]',null,'12207b06510193276b5fd9ad2fc55dcc004ada557d9259ca3505478bfef0b16ed977'),
  ('https://i.imgur.com/LJf4LzX.jpg', '2017-03-15 17:36:40', '2017-03-21 22:25:21','2017-03-21 22:25:20.88','2017-03-21 22:25:20.88384', 200,'image/jpeg', 'image/jpeg', 0, '','Puppies!','fe6d9fbd-32fe-4cf3-b48f-8f5010207f4c',0,0, '["Date","Tue, 21 Mar 2017 22:25:20 GMT"]',null,''),
  ('https://i.imgur.com/UE4nxKJ.gifv', '2017-03-15 17:36:40', '2017-03-21 22:25:21','2017-03-21 22:25:20.88','2017-03-21 22:25:20.88384', 200,'video/mp4', 'video/mp4', 0, '','Puppies!','e1a4eaff-1faf-48ea-9c2c-31968abc82bc',0,0, '["Date","Tue, 21 Mar 2017 22:25:20 GMT"]',null,''),
  ('https://i.imgur.com/ku6IEJf.gifv', '2017-03-15 17:36:40', '2017-03-21 22:25:21','2017-03-21 22:25:20.88','2017-03-21 22:25:20.88384', 200,'video/mp4', 'video/mp4', 0, '','Puppies!','8596e6b9-9bf6-45d6-b0c2-06a0f71de2df',0,0, '["Date","Tue, 21 Mar 2017 22:25:20 GMT"]',null,''),
  ('https://i.imgur.com/y22NjCp.jpg', '2017-03-15 17:36:40', '2017-03-21 22:25:21','2017-03-21 22:25:20.88','2017-03-21 22:25:20.88384', 200,'image/jpeg', 'image/jpeg', 0, '','Puppies!','98179ab7-8cd9-4d05-a6c8-24df846b8dd2',0,0, '["Date","Tue, 21 Mar 2017 22:25:20 GMT"]',null,''),
  ('https://i.imgur.com/1x8lR0p.jpg', '2017-03-15 17:36:40', '2017-03-21 22:25:21','2017-03-21 22:25:20.88','2017-03-21 22:25:20.88384', 200,'image/jpeg', 'image/jpeg', 0, '','Puppies!','41471973-9595-470f-b299-30a3423e267e',0,0, '["Date","Tue, 21 Mar 2017 22:25:20 GMT"]',null,''),
  ('http://i.imgur.com/26fVJAE.gifv', '2017-03-15 17:36:40', '2017-03-21 22:25:21','2017-03-21 22:25:20.88','2017-03-21 22:25:20.88384', 200,'video/mp4', 'video/mp4', 0, '','Puppies!','53d4c4bd-9802-41ad-82ae-1e82734d56fc',0,0, '["Date","Tue, 21 Mar 2017 22:25:20 GMT"]',null,'');
-- name: delete-urls
delete from urls;

-- name: insert-links
-- insert into links values
-- ('2017-01-01 00:00:02','2017-01-01 00:00:02','http://www.epa.gov','http://www.epa.gov');
-- name: delete-links
delete from links;

-- name: insert-snapshots
-- insert into snapshots values
--  ();
-- name: delete-snapshots
delete from snapshots;

-- name: insert-metadata
insert into metadata values
  -- hash,time_stamp,key_id,subject,prev,meta,deleted
  ('1220b499c5e883f2a3d47fe96c51779d05d5758c53c24caabc447db137b220c1688a', '2017-03-15 17:48:33', 'a5d6f8d8cbb15c0f60159b7fd58b1c555dbde076388ffce96c53f04b91ed377a', '12207b06510193276b5fd9ad2fc55dcc004ada557d9259ca3505478bfef0b16ed977', '', '{"description":"","title":"NAICS Classification Codes"}', false);
-- name: delete-metadata
delete from metadata;

-- name: insert-collections
insert into collections 
  (id, created, updated, 
    creator,
    title, description, url)
values
  ('6995febc-b7be-49ba-8297-1db68a703c3c','2017-07-13 15:57:32','2017-07-13 21:39:38',
    'EDGI_644b51b9567d0d999e40f697d7406a26030cde95a83775d285ff1f57a73b3ebc',
    'EPA TRU Datasets', 'essential TRU datasets', ''),
  ('f444f782-2110-43ca-956c-2c5f0dd56b1a','2017-07-12 23:18:59','2017-07-13 21:57:11',
    'EDGI_644b51b9567d0d999e40f697d7406a26030cde95a83775d285ff1f57a73b3ebc',
    'NOAA Volatile Organic Compound CSV files', 'VOC datasests', ''),
  ('a73a9d04-0fdb-40c8-a97f-288af36e8f6f','2017-07-13 18:45:03','2017-07-13 21:59:25',
    'blackglade_644b51b9567d0d999e40f697d7406a26030cde95a83775d285ff1f57a73b3ebc',
    'Test Collection', '', ''),
  ('4ed29ffe-150f-42ce-b0a9-2aadde646bcb','2017-07-11 17:50:19','2017-07-14 00:30:10',
    'jeffliu_644b51b9567d0d999e40f697d7406a26030cde95a83775d285ff1f57a73b3ebc',
    'All of the puppies', 'My fav puppy images', '');
-- name: delete-collections
delete from collections;

-- name: insert-collection_items
INSERT INTO collection_items
  (collection_id, url_id, index, description)
VALUES
  ('4ed29ffe-150f-42ce-b0a9-2aadde646bcb', 'fe6d9fbd-32fe-4cf3-b48f-8f5010207f4c', 0, 'puppies in a car'),
  ('4ed29ffe-150f-42ce-b0a9-2aadde646bcb', 'e1a4eaff-1faf-48ea-9c2c-31968abc82bc', 1, 'puppy in a sweater in a car'),
  ('4ed29ffe-150f-42ce-b0a9-2aadde646bcb', '8596e6b9-9bf6-45d6-b0c2-06a0f71de2df', 2, 'momma with puppies'),
  ('4ed29ffe-150f-42ce-b0a9-2aadde646bcb', '98179ab7-8cd9-4d05-a6c8-24df846b8dd2', 3, 'brown puppy'),
  ('4ed29ffe-150f-42ce-b0a9-2aadde646bcb', '41471973-9595-470f-b299-30a3423e267e', 4, ''),
  ('4ed29ffe-150f-42ce-b0a9-2aadde646bcb', '53d4c4bd-9802-41ad-82ae-1e82734d56fc', 5, '');
-- name: delete-collection_items
DELETE FROM collection_items;

-- name: insert-uncrawlables
insert into uncrawlables 
  ( id,url,created,updated,creator_key_id,
    name,email,event_name,agency_name,
    agency_id,subagency_id,org_id,suborg_id,subprimer_id,
    ftp,database,interactive,many_files,comments) 
values
  ( '55dd07ac-54cb-4f9d-b0a6-77d3d55c0d9e', 'https://www.census.gov/topics/economy/classification-codes.html', '2017-01-11 20:19:27', '2017-01-11 20:19:27','',
    'Sonal Ranjit','sonal.ranjit3@gmail.com','Testing','DJ Marvel Agency',
    '2','1','2','','234',
    false,false,false,false,'Advisory Committee on the Records of Congress');
-- name: delete-uncrawlables
delete from uncrawlables;

-- name: insert-archive_requests
-- insert into archive_requests values
--  ('8b14f3d6-882f-4dd5-92f8-abaac220864f','2017-01-01 00:00:01','http://www.apple.com','');
-- name: delete-archive_requests
delete from archive_requests;

-- name: insert-data_repos
insert into data_repos
  (id,created,updated,title,description,url)
values
  ('d5d9e72c-88f7-46ea-96cf-ae0f590a6f29', '2017-01-01 00:00:01', '2017-01-01 00:00:01', 'archivers 2.0', '', 'https://alpha.archivers.space'),
  ('2e410539-c8ae-4b58-9f46-7f56f954d502', '2017-01-01 00:00:01', '2017-01-01 00:00:01', 'archivers.space', '', 'https://www.archivers.space'),
  ('0268b749-f797-4c7e-bc88-cdb910ba4e6b', '2017-01-01 00:00:01', '2017-01-01 00:00:01', 'EOT Nomination Tool', '', 'https://github.com/edgi-govdata-archiving/eot-nomination-tool'),
  ('e7e78c62-4ef8-45ec-8373-77eedbd44b65', '2017-01-01 00:00:01', '2017-01-01 00:00:01', 'Internet Archive', '', 'https://archive.org'),
  ('0ed4b297-2af3-47f8-a746-0df780e0ea33', '2017-01-01 00:00:01', '2017-01-01 00:00:01', 'Project Svalbard', '', '');
-- name: delete-data_repos
delete from data_repos;