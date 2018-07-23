package application

import (
	"github.com/coreos/go-semver/semver"
	"github.com/jinzhu/gorm"
	"github.com/ovh/lhasa/api/db"
	"github.com/ovh/lhasa/api/hateoas"
	"github.com/ovh/lhasa/api/v1"
	"github.com/sirupsen/logrus"
)

// LatestUpdater updates the application latest version to the given version
type LatestUpdater func(*v1.Release, logrus.FieldLogger) error

// NewLatestUpdater instantiates a LatestUpdater
func NewLatestUpdater(tm db.TransactionManager, appRepoFactory RepositoryFactory, latestRepoFactory RepositoryLatestFactory) LatestUpdater {
	return func(version *v1.Release, log logrus.FieldLogger) error {
		return tm.Transaction(func(db *gorm.DB) error {
			appRepo := appRepoFactory(db)
			latestRepo := latestRepoFactory(db)

			log := log.WithFields(logrus.Fields{
				"domain":  version.Domain,
				"name":    version.Name,
				"version": version.Version,
			})
			application, err := latestRepo.FindApplication(version.Domain, version.Name)
			if err != nil && !hateoas.IsEntityDoesNotExistError(err) {
				return err
			}
			if hateoas.IsEntityDoesNotExistError(err) || application.LatestReleaseID == nil {
				log.Debug("application doesn't exist or doesn't have a latest so it will be created to the given version")
				application = &v1.Application{
					ID:     application.ID,
					Domain: version.Domain,
					Name:   version.Name,
				}
			}

			if err := appRepo.Save(version); err != nil {
				return err
			}
			if shouldUpdate(application.LatestRelease, version, log) {
				application.LatestRelease = version
				return latestRepo.Save(application)
			}
			return nil
		}, log)
	}
}

func shouldUpdate(current, submitted *v1.Release, log *logrus.Entry) bool {
	if current == nil {
		return true
	}
	if submitted == nil {
		return false
	}
	submittedSemver, err := semver.NewVersion(submitted.Version)
	if err != nil {
		log.Infof("version '%s' is not semver compliant, so it wont be used as latest", submitted.Version)
		return false
	}
	currentSemver, err := semver.NewVersion(current.Version)
	if err != nil {
		return true
	}
	return currentSemver.Compare(*submittedSemver) <= 0
}
