package main

import (
	"flag"
	"gopkg.in/mgo.v2"
	"time"
	"log"
)

var MongoAddress = flag.String("mongoaddr", "localhost", "the address to the mongo server <ip>:<port>")
var MongoTimeout = flag.Duration("mongotimeout", time.Second*10, "the timeout after the dial is cancelled")
var MongoUsername = flag.String("mongouser", "", "the username to authenticate with")
var MongoPassword = flag.String("mongopass", "", "the password to authenticate with")
var MongoDatabase = flag.String("mongoauthdb", "", "the database to authenticate with")
var MongoStorageDatabase = flag.String("mongostoragedb", "bandwidthmonitor", "the database to save the data in")
var MongoStorageCollection = flag.String("mongostoragecollection", "bandwidthentries", "the collection to save the data in")

var MonitoringDelay = flag.Duration("period", time.Minute*5, "the period in which bandwidth is monitored/measured")

var speedTester SpeedTester

func main() {
	speedTester = &SpeedTesterCLI{}
	flag.Parse()
	dialInfo := &mgo.DialInfo{
		Addrs:    []string{*MongoAddress},
		Timeout:  *MongoTimeout,
		Username: *MongoUsername,
		Password: *MongoPassword,
		Database: *MongoDatabase,
	}
	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		panic(err)
	}
	storageCollection := session.DB(*MongoStorageDatabase).C(*MongoStorageCollection)
	for {
		go func() {
			log.Println("Beginning new bandwidth monitor process...")
			old := time.Now()
			result, err := speedTester.MeasureBandwidth()
			if err != nil {
				log.Printf("There was an error while measuring the bandwidth: %v", err)
				return
			}
			if err = storageCollection.Insert(result); err != nil {
				log.Printf("There was an error while inserting the result: %v", err)
				return
			}
			log.Printf("Monitored bandwidth in %v seconds.", time.Since(old).Seconds())
		}()
		time.Sleep(*MonitoringDelay)
	}
}
