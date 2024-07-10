package pro

import (
	"github.com/nguyentrunghieu15/be-beehome-prj/authorize-service/internal/model"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/mongox"
)

func SetupRepo(client *mongox.ClientWrapper) {
	userRepository = mongox.Repository[model.User]{
		Client:     client,
		Collection: "users",
	}
	hireRepository = mongox.Repository[model.Hire]{
		Client:     client,
		Collection: "hires",
	}
	providerRepository = mongox.Repository[model.Provider]{
		Client:     client,
		Collection: "providers",
	}
	groupServiceRepository = mongox.Repository[model.GroupService]{
		Client:     client,
		Collection: "group_services",
	}
	serviceRepository = mongox.Repository[model.Service]{
		Client:     client,
		Collection: "services",
	}
	reviewRepository = mongox.Repository[model.Review]{Client: client, Collection: "reviews"}
}
