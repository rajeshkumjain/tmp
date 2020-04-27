package user

import (
	"context"
	"database/sql"
	"errors"
	"github.com/ecash/domain/entity"
	"github.com/ecash/domain/infra"
	repo "github.com/ecash/domain/repository"
	"log"
	"strings"
)

// InsertUser : Save Data in Entity RegisteredUser
func InsertUser(u *entity.RegisteredUser) (int64, error) {

	log.Println("*** Inside InsertUser Function ***")

	db, err := infra.NewDB()
	ctx := context.Background()
	if err != nil {
//		log.Println(repo.MessageMap["C001"])
		return 0, err
	}
	defer db.Close()

	emFound, _ := FindBy("email", strings.ToLower(strings.TrimSpace(u.Email)))
	mFound, _  := FindBy("mobile", u.Mobile)
	if emFound {
//		log.Println(repo.MessageMap["V001"])
		return 0, errors.New(repo.MessageMap["V001"])
	}
	if mFound {
		log.Println(repo.MessageMap["V002"])
		return 0, errors.New(repo.MessageMap["V002"])
	}

	tsql := `INSERT INTO RegisteredUser ( ru_firstname, ru_middlename, ru_surname, ru_email, ru_mobile, ru_password, ru_emailVerificationCode, 
		ru_mobileVerificationCode, ru_i_agree_t_n_c, ru_i_agree_newsletter, ru_registration_source, 
		ru_third_party_signin_id, ru_third_party_signin_name, ru_enable_flag, cf_active_flag, 
		cf_date_created_on, cf_date_last_updated_on, ru_is_mobile_verified, ru_is_email_verified ) VALUES (@Firstname, @Middlename, @Surname, @Email, @Mobile, @Password, @EmailCode,
			@MobileCode, @AgreeTermsCondition, @AgreeSendNewsletters, @RegistrationSource,
			@SocialPluginID, @SocialLoginPlugin, @EnableFlag, @ActiveFlag,
			@DateCreatedOn, @DateLastUpdatedOn, @MobileVerified, @EmailVerified  );SELECT convert(bigint, SCOPE_IDENTITY())`

	stmt, err := db.Prepare(tsql)
	defer stmt.Close()
	row := stmt.QueryRowContext(
		ctx,
		sql.Named("Firstname", u.Firstname),
		sql.Named("Surname", u.Surname),
		sql.Named("Middlename", u.Middlename),
		sql.Named("Email", strings.ToLower(strings.TrimSpace(u.Email))),
		sql.Named("Mobile", u.Mobile),
		sql.Named("Password", u.Password),
		sql.Named("EmailCode", u.EmailCode),
		sql.Named("MobileCode", u.MobileCode),
		sql.Named("AgreeTermsCondition", u.AgreeTermsCondition),
		sql.Named("AgreeSendNewsletters", u.AgreeSendNewsletters),
		sql.Named("RegistrationSource", u.RegistrationSource),
		sql.Named("SocialPluginID", u.SocialPluginID),
		sql.Named("SocialLoginPlugin", u.SocialLoginPlugin),
		sql.Named("EnableFlag", u.EnableFlag),
		sql.Named("ActiveFlag", u.ActiveFlag),
		sql.Named("DateCreatedOn", u.DateCreatedOn),
		sql.Named("DateLastUpdatedOn", u.DateLastUpdatedOn),
		sql.Named("MobileVerified", u.MobilVerified),
		sql.Named("EmailVerified", u.EmailVerified),
		)

	var userid int64
	err = row.Scan(&userid)
	if err != nil {
		//log.Println(repo.MessageMap["C002"], err.Error())
		return 0, errors.New( repo.MessageMap["C002"] + " " + err.Error())
	}
	log.Println("New User ID : ", userid)
	return userid, nil
}
// End of Insert Functionality


// FindBy : find by email / mobile in RegisteredUser table
func FindBy(findby string, s string) (bool, error) {
	var ErrorMessage string
	//log.Println("findby :", findby, " string to find :",s)
	if s =="" {
		ErrorMessage = "Search String Can't be blank"
		log.Println(ErrorMessage)
		return true, nil
	}
	db, err := infra.NewDB()
	ctx := context.Background()
	if err != nil {
		//log.Println(repo.MessageMap["C001"])
		return false, errors.New( repo.MessageMap["C001"] )
	}
	defer db.Close()
	var rid int64
	if findby == "email" {
		tsql := "SELECT ru_id FROM RegisteredUser where ru_email = @Email"
		rows, _ := db.QueryContext(ctx, tsql, sql.Named("Email", s))
		if rows.Next() {
			rows.Scan(&rid)
			//log.Println("Found Record : ", rid)
		}
		ErrorMessage = repo.MessageMap["V001"]
		defer rows.Close()
	} else {
		if findby == "mobile" {
			tsql := "SELECT ru_id FROM RegisteredUser where ru_mobile = @Mobile"
			rows, _ := db.QueryContext(ctx, tsql, sql.Named("Mobile", s))
			if rows.Next() {
				rows.Scan(&rid)
			//	log.Println("Found Record : ", rid)
			}
			ErrorMessage = repo.MessageMap["V002"]
			defer rows.Close()
		}
	}

	if rid > 0 {
		log.Println("ROW Found with ID", rid, ErrorMessage)
		return true, nil
	}
	return false, nil
}
