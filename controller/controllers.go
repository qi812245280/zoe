package controller

var (
	OrgController       *OrgControllerEngine
	UserController      *UserControllerEngine
	PrivilegeController *PrivilegeControllerEngine
	ProjectController   *ProjectControllerEngine
)

func Initialize() error {
	if e, err := NewOrgControllerEngine(); err != nil {
		return err
	} else {
		OrgController = e
	}

	if e, err := NewUserControllerEngine(); err != nil {
		return err
	} else {
		UserController = e
	}

	if e, err := NewPrivilegeControllerEngine(); err != nil {
		return err
	} else {
		PrivilegeController = e
	}

	if e, err := NewProjectControllerEngine(); err != nil {
		return err
	} else {
		ProjectController = e
	}

	return nil
}
