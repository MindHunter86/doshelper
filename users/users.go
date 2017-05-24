package users

type Users struct {
	e error
}
func (self *Users) InitModule() ( *Users, error ) {
	return nil,nil
}
func (self *Users) DeInitModule() error {
	return nil
}
