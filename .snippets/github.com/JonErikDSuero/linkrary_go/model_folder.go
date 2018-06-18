package main

import (
  "time"
  "gopkg.in/mgo.v2/bson"
)

type (
  Folders []Folder
  Folder struct {
    Id bson.ObjectId `json:"id" bson:"_id"`
    Name string `json:"name" bson:"name"`
    UpdatedAt time.Time `json:"updated_at" bsom:"created_at"`
  }
)


func (mc MgoCon) Folder_Upsert(folder *Folder) (err error) {
  if folder.Id.Hex() == "" {
    folder.Id = bson.NewObjectId()
  }
  folder.UpdatedAt = time.Now()
  _, err = mc.DB.C("folder").UpsertId(folder.Id, folder)
  return
}


func (mc MgoCon) Folder_All(folders *Folders) (err error) {
  err = mc.DB.C("folder").Find(bson.M{}).All(folders)
  return
}


func (mc MgoCon) Folder_Find(folder *Folder, query interface{}) (err error) {
  err = mc.DB.C("folder").Find(query).One(&folder)
  return
}

func (mc MgoCon) Folder_Suggest(folder_suggested *Folder, tags_filtered *[]string, extra_info *string) (err error) {
  //TODO: what if there are no links at all?
  var links Links
  var best_folder_id bson.ObjectId
  best_score := -1 // if all scores are 0's, the first folder will be default

  *tags_filtered = Tag_Filter(extra_info) // remove empty strings

  if err = mc.Link_All(&links); err != nil {
    panic(err)
  }

  scores := make(map[bson.ObjectId]int)
  for _, link := range links {
    scores[link.FolderId] += Tag_CommonalityScore(link.Tags, *tags_filtered)
  }

  for folder_id, score := range scores {
    if (best_score < score) {
      best_score = score
      best_folder_id = folder_id
    }
  }
  mc.DB.C("folder").Find(bson.M{"_id": best_folder_id}).One(folder_suggested)
  return err
}
