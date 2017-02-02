package auth

import (
	"strings"

	"github.com/qrclabs/nanogit/log"
	"github.com/qrclabs/nanogit/settings"
)

func CheckAuth(key string, path string) (read bool, write bool) {
	log.Trace("auth: CheckAuth, path: %s", path)
	org, repo := splitPath(cleanPath(path))
	orgRead, orgWrite := authOrg(key, org)
	repoRead, repoWrite := authRepo(key, repo)
	return orgRead || repoRead, orgWrite || repoWrite
}

func cleanPath(path string) string {
	return strings.Replace(path, "'", "", -1)
}

func splitPath(path string) (org string, repo string) {
	sliceStr := strings.Split(path, "/")
	return strings.ToLower(sliceStr[0]), strings.ToLower(sliceStr[1])
}

func authOrg(key string, orgPath string) (read bool, write bool) {
	log.Trace("auth: authOrg, org: %s", orgPath)
	orgConfig, err := settings.ConfInfo.LookupOrgById(orgPath)
	if err != nil {
		log.Error("auth: %v", err)
		return false, false
	}
	userConfig, err := settings.ConfInfo.LookupUserByKey(key)
	if err != nil {
		log.Error("auth: %v", err)
		return false, false
	}

	log.Trace("auth: authOrg, orgConfig.Id: %s", orgConfig.Id)
	// Loop on user orgs to compare teams access policy
	for _, userOrg := range userConfig.Orgs {
		// Only interested by org given in parameter
		if userOrg.Id == orgConfig.Id {
			// Loop on user teams
			for _, userTeam := range userOrg.Teams {
				// Loop on org teams
				for _, orgTeam := range orgConfig.Teams {
					// If we find a corresponding team,
					// return write and read policy
					if userTeam == orgTeam.Name {
						return orgTeam.Read, orgTeam.Write
					}
				}
			}
		}
	}
	return false, false
}

func authRepo(key string, repoPath string) (read bool, write bool) {
	log.Trace("auth: authRepo, repo: %s", repoPath)
	return true, true
}
