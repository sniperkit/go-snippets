/*
* @Author: souravray
* @Date:   2015-01-24 10:34:31
* @Last Modified by:   souravray
* @Last Modified time: 2015-01-25 23:34:32
 */

package controllers

import (
	"fmt"
	"github.com/gophergala/tinyembassy/site/models"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"net/http"
)

func CreateCampaignPage(rw http.ResponseWriter, req *http.Request) {
	if !IsAuth(req) {
		http.Redirect(rw, req, "/", 301)
	}

	websession, _ := store.Get(req, "pp-session")
	advertiser := websession.Values["id"].(*models.Advertiser)
	if root, ok := hasCampaign(advertiser.Id); ok {
		redirctpath := fmt.Sprint("/b/", root)
		http.Redirect(rw, req, redirctpath, 301)
	}

	render(rw, "campaignregistration.html")
	return
}

func CreateCampaign(rw http.ResponseWriter, req *http.Request) {
	if !IsAuth(req) {
		http.Redirect(rw, req, "/", 301)
	}

	websession, _ := store.Get(req, "pp-session")
	advertiser := websession.Values["id"].(*models.Advertiser)

	if root, ok := hasCampaign(advertiser.Id); ok {
		redirctpath := fmt.Sprint("/b/", root)
		http.Redirect(rw, req, redirctpath, 301)
	}
	rootURL := req.FormValue("urlroot")

	session, err := mgo.Dial(conf.DbURI)
	if err != nil {
		panic(err)
	}
	c := session.DB(conf.DbName).C("campaigns")
	defer session.Close()
	result := models.Campaign{}

	err = c.Find(bson.M{"url_root": rootURL}).One(&result)
	if err != nil {
		doc := models.Campaign{Id: bson.NewObjectId(), AdvertiserId: advertiser.Id, URLRoot: rootURL}
		err = c.Insert(doc)
		if err != nil {
			fmt.Printf("Can't insert document: %v\n", err)
			http.Redirect(rw, req, "/campaign/create", 301)
		} else {
			redirctpath := fmt.Sprint("/b/", rootURL)
			http.Redirect(rw, req, redirctpath, 301)
		}
	} else {
		fmt.Println("URL already exist..")
		http.Redirect(rw, req, "/campaign/create", 301)
	}
	session.Close()
	return
}

func hasCampaign(advertiserId bson.ObjectId) (string, bool) {
	session, err := mgo.Dial(conf.DbURI)
	if err != nil {
		panic(err)
	}
	c := session.DB(conf.DbName).C("campaigns")
	defer session.Close()
	result := models.Campaign{}

	err = c.Find(bson.M{"advertiser_id": advertiserId}).One(&result)
	if err == nil {
		return result.URLRoot, true
	}
	return "", false
}

func CBG(rw http.ResponseWriter, req *http.Request) {
	render(rw, "CBG.html")
	return
}

// func CreateBadgeGroup(rw http.ResponseWriter, req *http.Request) {
// 	fmt.Println("in Create Badge Group....")
// 	targetURL := req.FormValue("targetURL")
// 	title := req.FormValue("title")
// 	campaignName := req.FormValue("campaignName")

// 	websession, _ := store.Get(req, "pp-session")
// 	fmt.Println(websession)

// 	advertiser := websession.Values["id"].(*models.Advertiser)
// 	fmt.Println(websession)
// 	session, err := mgo.Dial(conf.DbURI)
// 	if err != nil {
// 		panic(err)
// 	}
// 	c := session.DB(conf.DbName).C("campaigns")

// 	defer session.Close()

// 	campaign := models.Campaign{}
// 	fmt.Println("Search for " + campaignName)
// 	err = c.Find(bson.M{"campaignName": campaignName, "advertiser_id": advertiser.Id}).One(&campaign)
// 	if err == nil {
// 		bc := session.DB(conf.DbName).C("badgeGroup")
// 		badgeGroup := models.BadgeGroup{}
// 		err = bc.Find(bson.M{"title": title, "campaign_id": campaign.CampaignId}).One(&badgeGroup)
// 		if err != nil {
// 			doc := models.BadgeGroup{BadgeGroupId: bson.NewObjectId(), CampaignId: campaign.CampaignId, Title: title, TargetURL: targetURL}
// 			err = bc.Insert(doc)
// 			if err != nil {
// 				fmt.Printf("Can't insert document: %v\n", err)
// 				http.Redirect(rw, req, "/group/create", http.StatusTemporaryRedirect)
// 			} else {
// 				http.Redirect(rw, req, "/badge/create", http.StatusTemporaryRedirect)
// 			}
// 		} else {
// 			fmt.Println("badge already exist...")
// 			http.Redirect(rw, req, "/group/create", http.StatusTemporaryRedirect)
// 		}
// 	} else {
// 		fmt.Println("Campaign does not exit")
// 		fmt.Println(err)
// 		http.Redirect(rw, req, "/group/create", http.StatusTemporaryRedirect)
// 	}
// 	session.Close()
// 	return
// }

// func GetCampaignData(rw http.ResponseWriter, req *http.Request) {
// 	fmt.Println("GetCampaignData....")
// 	//TODO: return campaign data for given Advertiser
// }
