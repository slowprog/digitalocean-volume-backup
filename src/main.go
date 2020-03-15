package main

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/digitalocean/godo"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"

	"github.com/slowprog/digitalocean-volumes-backup/src/settings"
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	log.Info("Volumes backup begin")

	conf := settings.NewConfig()

	oauthClient := oauth2.NewClient(context.Background(), settings.NewTokenSource(conf.AccessToken))

	client := godo.NewClient(oauthClient)

	ctx := context.TODO()

	volumes, response, err := client.Storage.ListVolumes(ctx, &godo.ListVolumeParams{})

	if err != nil {
		log.WithFields(log.Fields{
			"message":    err,
			"rate_limit": response,
		}).Fatal("Something bad happened with ListVolumes")
	}

	volumes = filterVolumesByName(volumes, conf.VolumesBackup)

	currentDate := time.Now()

	for _, vol := range volumes {
		snapshotName := fmt.Sprintf("%s-%s-%s", conf.SnapshotsPrefix, vol.Name, currentDate.Format("2006.01.02.15.04.05"))

		log.WithFields(log.Fields{
			"volume_id":     vol.ID,
			"volume_name":   vol.Name,
			"snapshot_name": snapshotName,
		}).Info("Creating snapshot")

		snapshot, response, err := client.Storage.CreateSnapshot(ctx, &godo.SnapshotCreateRequest{
			VolumeID: vol.ID,
			Name:     snapshotName,
		})

		if err != nil {
			log.WithFields(log.Fields{
				"message":    err,
				"rate_limit": response,
			}).Error("Something bad happened with CreateSnapshot")

			continue
		}

		log.WithFields(log.Fields{
			"volume_id":     vol.ID,
			"volume_name":   vol.Name,
			"snapshot_id":   snapshot.ID,
			"snapshot_name": snapshot.Name,
		}).Info("Snapshot has been created")

		snapshots, response, err := client.Storage.ListSnapshots(ctx, vol.ID, &godo.ListOptions{})

		if err != nil {
			log.WithFields(log.Fields{
				"message":    err,
				"rate_limit": response,
			}).Error("Something bad happened with ListSnapshots")

			continue
		}

		snapshots = filterSnapshotsByPrefix(snapshots, conf.SnapshotsPrefix)

		if len(snapshots) > conf.SnapshotsMax {
			sort.Slice(snapshots, func(i, j int) bool {
				date1, _ := time.Parse(time.RFC3339, snapshots[i].Created)
				date2, _ := time.Parse(time.RFC3339, snapshots[j].Created)

				return date1.Before(date2)
			})

			snapshots = snapshots[:len(snapshots)-conf.SnapshotsMax]

			log.WithFields(log.Fields{
				"snapshots_names": getSpapchotsNames(snapshots),
			}).Info("Removing old snapshots")

			for _, snapshot := range snapshots {
				response, err := client.Storage.DeleteSnapshot(ctx, snapshot.ID)

				if err != nil {
					log.WithFields(log.Fields{
						"message":    err,
						"rate_limit": response,
					}).Error("Something bad happened with DeleteSnapshot")

					continue
				}
			}
		}
	}

	log.Info("Volumes backup is over")
}

func filterSnapshotsByPrefix(snapshots []godo.Snapshot, snapshotPrefix string) []godo.Snapshot {
	log.WithFields(log.Fields{
		"snapshots_number": len(snapshots),
		"snapshots_names":  getSpapchotsNames(snapshots),
	}).Info("Exists snapshots of volume")

	result := make([]godo.Snapshot, 0)

	for _, snapshot := range snapshots {
		if strings.Index(snapshot.Name, snapshotPrefix) == 0 {
			result = append(result, snapshot)
		}
	}

	log.WithFields(log.Fields{
		"snapshots_number": len(snapshots),
		"snapshots_names":  getSpapchotsNames(snapshots),
		"snapshot_prefix":  snapshotPrefix,
	}).Info("Filtered snapshots of volume by prefix")

	return result
}

func getVolumesNames(volumes []godo.Volume) []string {
	result := make([]string, 0)

	for _, volume := range volumes {
		result = append(result, volume.Name)
	}

	return result
}

func getSpapchotsNames(snapshots []godo.Snapshot) []string {
	result := make([]string, 0)

	for _, snapshot := range snapshots {
		result = append(result, snapshot.Name)
	}

	return result
}

func filterVolumesByName(volumes []godo.Volume, volumesBackup []string) []godo.Volume {
	log.WithFields(log.Fields{
		"volumes_number": len(volumes),
		"volumes_names":  getVolumesNames(volumes),
	}).Info("Exists volumes")

	if len(volumesBackup) < 1 {
		return volumes
	}

	result := make([]godo.Volume, 0)

	InVolumesNeeded := func(name string) bool {
		for _, v := range volumesBackup {
			if v == name {
				return true
			}
		}

		return false
	}

	for _, volume := range volumes {
		if InVolumesNeeded(volume.Name) {
			result = append(result, volume)
		}
	}

	log.WithFields(log.Fields{
		"volumes_number": len(result),
		"volumes_names":  getVolumesNames(result),
		"volumes_needed": volumesBackup,
	}).Info("Filtered volumes")

	return result
}
