package dbc

import (
	"gorm.io/gorm"
	"time"

	"github.com/vault-thirteen/BytePackedPassword"
	"github.com/vault-thirteen/TR1/src/models/common"
)

const (
	CountOnError = -1
)

// Common methods.

func (dbc *DbController) countAllItems(model *gorm.DB) (n int, err error) {
	var n64 int64
	tx := model.Count(&n64)
	if tx.Error != nil {
		return CountOnError, tx.Error
	}
	return int(n64), nil
}
func (dbc *DbController) listItemsOnPage(model *gorm.DB, page int, dst any) (err error) {
	tx := model.Limit(dbc.pageSize).Offset((page - 1) * dbc.pageSize).Find(dst)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (dbc *DbController) listItemsOnPageWithTotalCount(model *gorm.DB, page int, dst any) (totalCount int, err error) {
	totalCount, err = dbc.countAllItems(model)
	if err != nil {
		return CountOnError, err
	}

	err = dbc.listItemsOnPage(model, page, dst)
	if err != nil {
		return CountOnError, err
	}

	return totalCount, nil
}

// User registration.

func (dbc *DbController) IsUserNameFree(userName string) (isFree bool, err error) {
	var n int64
	tx := dbc.db.Model(&cm.User{}).Where("name = ?", userName).Count(&n)
	if tx.Error != nil {
		return false, tx.Error
	}
	if n > 0 {
		return false, nil
	}
	return true, nil
}
func (dbc *DbController) IsUserEmailFree(userEmail string) (isFree bool, err error) {
	var n int64
	tx := dbc.db.Model(&cm.User{}).Where("email = ?", userEmail).Count(&n)
	if tx.Error != nil {
		return false, tx.Error
	}
	if n > 0 {
		return false, nil
	}
	return true, nil
}
func (dbc *DbController) ExistsRegistrationRequestWithUserName(userName string) (exists bool, err error) {
	var n int64
	tx := dbc.db.Model(&cm.RegistrationRequest{}).Where("user_name = ?", userName).Count(&n)
	if tx.Error != nil {
		return false, tx.Error
	}
	if n > 0 {
		return true, nil
	}
	return false, nil
}
func (dbc *DbController) ExistsRegistrationRequestWithUserEmail(userEmail string) (exists bool, err error) {
	var n int64
	tx := dbc.db.Model(&cm.RegistrationRequest{}).Where("user_email = ?", userEmail).Count(&n)
	if tx.Error != nil {
		return false, tx.Error
	}
	if n > 0 {
		return true, nil
	}
	return false, nil
}
func (dbc *DbController) CreateRegistrationRequest(rr cm.RegistrationRequest) (err error) {
	tx := dbc.db.Create(&rr)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (dbc *DbController) FindRegistrationRequestNRFA(rr *cm.RegistrationRequest) (err error) {
	tx := dbc.db.First(rr, "request_id = ? AND is_ready_for_approval = ?", rr.RequestId, false)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (dbc *DbController) MarkRegistrationRequestAsReadyForApproval(rr *cm.RegistrationRequest) (err error) {
	tx := dbc.db.Model(&rr).Update("is_ready_for_approval", true)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (dbc *DbController) CreateUser(user *cm.User, password string) (err error) {
	var pwd = &cm.Password{
		UserId: user.Id,
	}

	pwd.Bytes, err = bpp.PackSymbols([]rune(password))
	if err != nil {
		return err
	}

	user.Password = pwd

	tx := dbc.db.Create(user)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}
func (dbc *DbController) DeleteRegistrationRequestRFA(rr *cm.RegistrationRequest) (err error) {
	tx := dbc.db.Where("is_ready_for_approval = ?", true).Delete(&rr)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (dbc *DbController) DeleteRegistrationRequestNRFA(rr *cm.RegistrationRequest) (err error) {
	tx := dbc.db.Where("is_ready_for_approval = ?", false).Delete(&rr)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (dbc *DbController) GetFirstOutdatedRegistrationRequest(edgeTime time.Time) (rrs []cm.RegistrationRequest, err error) {
	tx := dbc.db.Limit(1).Where("is_ready_for_approval = ? AND created_at <= ?", false, edgeTime).Find(&rrs)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return rrs, nil
}
func (dbc *DbController) ListRegistrationRequestsRFA(page int) (rrs []cm.RegistrationRequest, totalCount int, err error) {
	model := dbc.db.Model(&cm.RegistrationRequest{}).Where("is_ready_for_approval = ?", true).Omit("UserPassword")

	totalCount, err = dbc.listItemsOnPageWithTotalCount(model, page, &rrs)
	if err != nil {
		return nil, totalCount, err
	}

	return rrs, totalCount, nil
}
func (dbc *DbController) GetRegistrationRequestRFA(userEmail string, rr *cm.RegistrationRequest) (err error) {
	tx := dbc.db.First(rr, "is_ready_for_approval = ? AND user_email = ?", true, userEmail)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

// User logging in.

func (dbc *DbController) ExistsLogInRequestWithUserEmail(userEmail string) (exists bool, err error) {
	var n int64
	tx := dbc.db.Model(&cm.LogInRequest{}).Where("user_email = ?", userEmail).Count(&n)
	if tx.Error != nil {
		return false, tx.Error
	}
	if n > 0 {
		return true, nil
	}
	return false, nil
}
func (dbc *DbController) GetUserByEmailAbleToLogIn(user *cm.User) (err error) {
	tx := dbc.db.First(user, "email = ? AND can_log_in = ?", user.Email, true)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (dbc *DbController) CreateLogInRequest(lir cm.LogInRequest) (err error) {
	tx := dbc.db.Create(&lir)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (dbc *DbController) GetFirstOutdatedLogInRequest(edgeTime time.Time) (lirs []cm.LogInRequest, err error) {
	tx := dbc.db.Limit(1).Where("created_at <= ?", edgeTime).Find(&lirs)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return lirs, nil
}
func (dbc *DbController) DeleteOldLogInRequest(lir *cm.LogInRequest) (err error) {
	tx := dbc.db.Delete(&lir)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (dbc *DbController) FindLogInRequest(lir *cm.LogInRequest) (err error) {
	tx := dbc.db.First(lir, "request_id = ?", lir.RequestId)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (dbc *DbController) GetUserByIdAbleToLogIn(user *cm.User) (err error) {
	tx := dbc.db.Preload("Password").First(user, "id = ? AND can_log_in = ?", user.Id, true)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (dbc *DbController) CreateSession(session *cm.Session) (err error) {
	tx := dbc.db.Create(session)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (dbc *DbController) DeleteLogInRequest(lir *cm.LogInRequest) (err error) {
	tx := dbc.db.Where("id = ?", lir.Id).Delete(&lir)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (dbc *DbController) CreateLogEvent(le *cm.LogEvent) (err error) {
	tx := dbc.db.Create(le)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

// User authorisation.

func (dbc *DbController) GetUserWithSessionByIdAbleToLogIn(user *cm.User) (err error) {
	tx := dbc.db.Preload("Session").First(user, "id = ? AND can_log_in = ?", user.Id, true)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (dbc *DbController) CreateLogOutRequest(lor cm.LogOutRequest) (err error) {
	tx := dbc.db.Create(&lor)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (dbc *DbController) GetFirstOutdatedLogOutRequest(edgeTime time.Time) (lors []cm.LogOutRequest, err error) {
	tx := dbc.db.Limit(1).Where("created_at <= ?", edgeTime).Find(&lors)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return lors, nil
}
func (dbc *DbController) DeleteOldLogOutRequest(lor *cm.LogOutRequest) (err error) {
	tx := dbc.db.Delete(&lor)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (dbc *DbController) FindLogOutRequest(lor *cm.LogOutRequest) (err error) {
	tx := dbc.db.First(lor, "request_id = ?", lor.RequestId)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (dbc *DbController) DeleteSession(session *cm.Session) (err error) {
	tx := dbc.db.Delete(session)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (dbc *DbController) DeleteLogOutRequest(lor *cm.LogOutRequest) (err error) {
	tx := dbc.db.Where("id = ?", lor.Id).Delete(&lor)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (dbc *DbController) GetFirstOutdatedSession(edgeTime time.Time) (ss []cm.Session, err error) {
	tx := dbc.db.Limit(1).Where("created_at <= ?", edgeTime).Find(&ss)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return ss, nil
}
func (dbc *DbController) DeleteOldSession(s *cm.Session) (err error) {
	tx := dbc.db.Delete(&s)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

// E-mail change.

func (dbc *DbController) CreateEmailChangeRequest(ecr cm.EmailChangeRequest) (err error) {
	tx := dbc.db.Create(&ecr)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (dbc *DbController) ExistsEmailChangeRequestWithNewEmail(newEmail string) (exists bool, err error) {
	var n int64
	tx := dbc.db.Model(&cm.EmailChangeRequest{}).Where("new_email = ?", newEmail).Count(&n)
	if tx.Error != nil {
		return false, tx.Error
	}
	if n > 0 {
		return true, nil
	}
	return false, nil
}
func (dbc *DbController) FindEmailChangeRequest(ecr *cm.EmailChangeRequest) (err error) {
	tx := dbc.db.First(ecr, "request_id = ?", ecr.RequestId)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (dbc *DbController) SaveUserEmail(user *cm.User, email string) (err error) {
	user.Email = email

	tx := dbc.db.Save(user)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (dbc *DbController) DeleteEmailChangeRequest(ecr *cm.EmailChangeRequest) (err error) {
	tx := dbc.db.Where("id = ?", ecr.Id).Delete(&ecr)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (dbc *DbController) GetFirstOutdatedEmailChangeRequest(edgeTime time.Time) (ecrs []cm.EmailChangeRequest, err error) {
	tx := dbc.db.Limit(1).Where("created_at <= ?", edgeTime).Find(&ecrs)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return ecrs, nil
}
func (dbc *DbController) DeleteOldEmailChangeRequest(ecr *cm.EmailChangeRequest) (err error) {
	tx := dbc.db.Delete(&ecr)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

// Password change.

func (dbc *DbController) ExistsPasswordChangeRequestWithUserId(user *cm.User) (exists bool, err error) {
	var n int64
	tx := dbc.db.Model(&cm.PasswordChangeRequest{}).Where("user_id = ?", user.Id).Count(&n)
	if tx.Error != nil {
		return false, tx.Error
	}
	if n > 0 {
		return true, nil
	}
	return false, nil
}
func (dbc *DbController) CreatePasswordChangeRequest(pcr cm.PasswordChangeRequest) (err error) {
	tx := dbc.db.Create(&pcr)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (dbc *DbController) GetFirstOutdatedPasswordChangeRequest(edgeTime time.Time) (pcrs []cm.PasswordChangeRequest, err error) {
	tx := dbc.db.Limit(1).Where("created_at <= ?", edgeTime).Find(&pcrs)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return pcrs, nil
}
func (dbc *DbController) DeleteOldPasswordChangeRequest(pcr *cm.PasswordChangeRequest) (err error) {
	tx := dbc.db.Delete(&pcr)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (dbc *DbController) FindPasswordChangeRequest(pcr *cm.PasswordChangeRequest) (err error) {
	tx := dbc.db.First(pcr, "request_id = ?", pcr.RequestId)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (dbc *DbController) SaveUserPassword(user *cm.User, password string) (err error) {
	user.Password.Bytes, err = bpp.PackSymbols([]rune(password))
	if err != nil {
		return err
	}

	//tx := dbc.db.Save(user) <- This does not work in modern versions of GORM !
	tx := dbc.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(user)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (dbc *DbController) DeletePasswordChangeRequest(pcr *cm.PasswordChangeRequest) (err error) {
	tx := dbc.db.Where("id = ?", pcr.Id).Delete(&pcr)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
