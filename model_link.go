package main

import (
  "time"
  "gopkg.in/mgo.v2/bson"
)


type (
  Links []Link
  Link struct {
    Id bson.ObjectId `json:"id" bson:"_id"`
    Name string `json:"name" bson:"name"`
    Url string `json:"url" bson:"url"`
    FolderId bson.ObjectId `json:"folder_id" bson:"folder_id"`
    Active bool `json:"active" bson:"active"`
    Tags []string `json:"tags" bson:"tags"`
    UpdatedAt time.Time `json:"updated_at" bsom:"created_at"`
  }
)


func (mc MgoCon) Link_Upsert(link *Link) (err error) {
  //Example: info, err := collection.Upsert(bson.M{"_id": id}, updated_object)
  if (link.Id.Hex() == "") {
    link.Id = bson.NewObjectId()
  }
  link.UpdatedAt = time.Now()
  _, err = mc.DB.C("link").Upsert(bson.M{"url": link.Url}, link)
  return
}


func (mc MgoCon) Link_All(links *Links) (err error) {
  err = mc.DB.C("link").Find(bson.M{}).All(links)
  return
}


func (mc MgoCon) Link_Find(link *Link, query interface{}) (err error) {
  err = mc.DB.C("link").Find(query).One(&link)
  return err
}

