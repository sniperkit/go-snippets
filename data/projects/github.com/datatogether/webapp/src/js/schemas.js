import { Schema, arrayOf } from 'normalizr';

// We use this Normalizr schemas to transform API responses from a nested form
// to a flat form where repos and users are placed in `entities`, and nested
// JSON objects are replaced with their IDs. This is very convenient for
// consumption by reducers, because we can easily build a normalized tree
// and keep it updated as we fetch more data.

// Read more about Normalizr: https://github.com/paularmstrong/normalizr
const sessionUserSchema = new Schema('session');
const keySchema = new Schema('keys', { idAttribute: 'sha256' });
const userSchema = new Schema('users');
const fileSchema = new Schema('files');
const groupsSchema = new Schema('groups');
const metadataSchema = new Schema('metadata', { idAttribute: metadata => `${metadata.keyId}.${metadata.subject}` });
const searchResultSchema = new Schema('searchResults', { idAttribute: 'url' });
const urlSchema = new Schema('urls', { idAttribute: 'url' });
const linkSchema = new Schema('links', { idAttribute: link => `${link.src.url}.${link.dst.url}` });
const consensusSchema = new Schema('consensus', { idAttribute: 'subject' });
const contentSchema = new Schema('content', { idAttribute: 'hash' });
const collectionSchema = new Schema('collections');
const primerSchema = new Schema('primers');
const sourceSchema = new Schema('sources');
const uncrawlableSchema = new Schema('uncrawlable');
const taskSchema = new Schema('tasks');
const collectionItemSchema = new Schema('collectionItems');

primerSchema.define({
  subPrimers: arrayOf(primerSchema),
});

fileSchema.new = (attrs) => {
  return Object.assign({}, attrs);
}

metadataSchema.new = (attrs) => {
  return Object.assign({}, attrs);
};

collectionSchema.new = (attrs) => {
  return Object.assign({ id: "new" }, attrs);
};

// Schemas for Github API responses.
const Schemas = {
  // currently-logged-in user
  SESSION_USER: sessionUserSchema,
  // users. many users. much parcipation.
  USER: userSchema,
  USER_ARRAY: arrayOf(userSchema),
  // user key management
  KEY: keySchema,
  KEY_ARRAY: arrayOf(keySchema),

  // groups of users, currently unimplemeneted
  GROUP: groupsSchema,
  GROUP_ARRAY: arrayOf(groupsSchema),

  PRIMER: primerSchema,
  PRIMER_ARRAY: arrayOf(primerSchema),

  SOURCE: sourceSchema,
  SOURCE_ARRAY: arrayOf(sourceSchema),
  // an external url for archiving
  URL: urlSchema,
  URL_ARRAY: arrayOf(urlSchema),
  // link from a src to a dst url
  LINK: linkSchema,
  LINK_ARRAY: arrayOf(linkSchema),

  // structured content, the basis of archiving. content has a shasum
  CONTENT: contentSchema,
  CONTENT_ARRAY: arrayOf(contentSchema),
  // metadata is a user-contributed object that describes content
  METADATA: metadataSchema,
  METADATA_ARRAY: arrayOf(metadataSchema),
  // consensus is the sum of all metadata for a given piece of content
  CONSENSUS: consensusSchema,
  CONSENSUS_ARRAY: arrayOf(consensusSchema),
  // collections are lists of content. they are also content themselves, and can have metadata
  // collections are "special content" on this webapp b/c users can edit them
  COLLECTION: collectionSchema,
  COLLECTION_ARRAY: arrayOf(collectionSchema),

  COLLECTION_ITEM: collectionItemSchema,
  COLLECTION_ITEM_ARRAY: arrayOf(collectionItemSchema),

  // search results can be a number of different models
  // searchResult wraps those discrete types
  SEARCH_RESULT: searchResultSchema,
  SEARCH_RESULT_ARRAY: arrayOf(searchResultSchema),

  TASK: taskSchema,
  TASK_ARRAY: arrayOf(taskSchema),

  // TODO - finish comment
  UNCRAWLABLE: uncrawlableSchema,
  UNCRAWLABLE_ARRAY: arrayOf(uncrawlableSchema),

  // TODO - finish comment
  FILE: fileSchema,
  FILE_ARRAY: arrayOf(fileSchema),
};

export default Schemas;
