package model

import (
	"database/sql"
	"time"

	"fmt"
)

type MetaChatRequest struct {
	CustomerId  int64  `json:"customer_id"`
	SiteId      int64  `json:"site_id"`
	WebsiteUuid string `json:"website_uuid"`
	Title       string `json:"title"`
	Email       string `json:"email"`
	Domain      string `json:"domain"`
}

type MetachatRow struct {
	MetachatId           int64
	MetachatCustomerId   int64
	MetachatSiteId       int64
	MetachatThemeId      int64
	MetachatIsEnabled    bool
	MetachatCreationDate time.Time
	MetachatIsDeleted    bool
	// Appearance fields
	AppearanceId                       int64
	AppearanceTitle                    string
	AppearanceLayout                   int64
	AppearanceLayoutHeaderColor        string
	AppearanceLayoutHeaderText         string
	AppearanceLayoutBackGround         string
	AppearanceLayoutFooter             string
	AppearanceLayoutOperatorBubble     string
	AppearanceLayoutOperatorBubbleText string
	AppearanceLayoutVisitorBubble      string
	AppearanceLayoutVisitorBubbleText  string
	AppearanceBtnPosition              int64
	AppearanceShowAvatar               bool
	AppearanceAvatarId                 string
	AppearanceIconButton               int64
	AppearanceShowPositionNum          bool
	AppearancePlaySoundStart           bool
	AppearancePlaySoundDecre           bool
	AppearanceOfferTalkBot             bool
	//AppearanceLogo                     []byte
	// Deployment fields
	DeploymentId                      int64
	DeploymentDeployHtmlCode          string
	DeploymentMailToSendHtmlCode      string
	DeploymentDomainAllowedDeployment string
	// prechat fields
	PrechatId               int64
	PrechatEnabled          bool
	PrechatName             bool
	PrechatNameRequired     bool
	PrechatLastname         bool
	PrechatLastnameRequired bool
	PrechatEmail            bool
	PrechatEmailRequired    bool
	PrechatQuestion         bool
	PrechatQuestionRequired bool
	PrechatCustom1          string
	PrechatCustom1Required  bool
	PrechatCustom2          string
	PrechatCustom2Required  bool
	PrechatCustom3          string
	PrechatCustom3Required  bool
	// Postchat fields
	PostchatId         int64
	PostchatEnabled    bool
	PostchatSendSurvey bool
	PostchatAsk        bool
	PostchatShowSolved bool
	PostchatSendEndMsg bool
}

// this model is used to get data from other microservices
type Site struct {
	Id           int64     `json:"id"`
	CustomerId   int64     `json:"customer_id"`
	Name         string    `json:"name"`
	Domain       string    `json:"domain"`
	Description  string    `json:"description"`
	IsEnabled    bool      `json:"is_enabled"`
	CreationDate time.Time `json:"creation_date"`
	IsDeleted    bool      `json:"is_deleted"`
}

type Metachat struct {
	Id              int64       `json:"id"`
	CustomerId      int64       `json:"customer_id"`
	SiteId          int64       `json:"site_id"`
	WebsiteUuid     string      `json:"website_uuid"`
	SiteTextDescrip string      `json:"site_text_descrip"`
	ThemeId         int64       `json:"theme_id"`
	Appearance      *Appearance `json:"appearance"`
	Deployment      *Deployment `json:"deployment"`
	Prechat         *Prechat    `json:"preChat"`
	Postchat        *Postchat   `json:"postChat"`
	IsEnabled       bool        `json:"is_enabled"`
	CreationDate    time.Time   `json:"creation_date"`
	IsDeleted       bool        `json:"is_deleted"`
}

type Appearance struct {
	Id                       int64  `json:"id"`
	Title                    string `json:"title"`
	Layout                   int64  `json:"layout"`
	LayoutHeaderColor        string `json:"layout_header_color"`
	LayoutHeaderText         string `json:"layout_header_text"`
	LayoutBackGround         string `json:"layout_background"`
	LayoutFooter             string `json:"layout_footer"`
	LayoutOperatorBubble     string `json:"layout_operator_bubble"`
	LayoutOperatorBubbleText string `json:"layout_operator_bubble_text"`
	LayoutVisitorBubble      string `json:"layout_visitor_bubble"`
	LayoutVisitorBubbleText  string `json:"layout_visitor_bubble_text"`
	BtnPosition              int64  `json:"btn_position"`
	ShowAvatar               bool   `json:"show_avatar"`
	AvatarId                 string `json:"avatar_id"`
	IconButton               int64  `json:"icon_button"`
	ShowPositionNum          bool   `json:"show_position_num"`
	PlaySoundStart           bool   `json:"play_sound_start"`
	PlaySoundDecre           bool   `json:"play_sound_decre"`
	OfferTalkBot             bool   `json:"offer_talk_bot"`
	//Logo                     []byte `json:"logo"`
}

type Deployment struct {
	Id                      int64  `json:"id"`
	DeployHtmlCode          string `json:"deploy_html_code"`
	MailToSendHtmlCode      string `json:"mail_to_send_html_code"`
	DomainAllowedDeployment string `json:"domain_allowed_deployment"`
}

type Prechat struct {
	Id               int64  `json:"id"`
	Enabled          bool   `json:"enabled"`
	Name             bool   `json:"name"`
	NameRequired     bool   `json:"name_required"`
	Lastname         bool   `json:"lastname"`
	LastnameRequired bool   `json:"lastname_required"`
	Email            bool   `json:"email"`
	EmailRequired    bool   `json:"email_required"`
	Question         bool   `json:"question"`
	QuestionRequired bool   `json:"question_required"`
	Custom1          string `json:"custom1"`
	Custom1Required  bool   `json:"custom1_required"`
	Custom2          string `json:"custom2"`
	Custom2Required  bool   `json:"custom2_required"`
	Custom3          string `json:"custom3"`
	Custom3Required  bool   `json:"custom3_required"`
}

type Postchat struct {
	Id         int64 `json:"id"`
	Enabled    bool  `json:"enabled"`
	SendSurvey bool  `json:"send_survey"`
	Ask        bool  `json:"ask"`
	ShowSolved bool  `json:"show_solved"`
	SendEndMsg bool  `json:"send_end_msg"`
}

type Metachats []*Metachat

func (metachat *Metachat) loadFromMetachat(s MetachatRow) {

	metachat.Id = s.MetachatId
	metachat.CustomerId = s.MetachatCustomerId
	metachat.SiteId = s.MetachatSiteId
	metachat.ThemeId = s.MetachatThemeId
	metachat.IsEnabled = s.MetachatIsEnabled
	metachat.CreationDate = s.MetachatCreationDate
	metachat.IsDeleted = s.MetachatIsDeleted

	metachat.Appearance = &Appearance{}
	metachat.Appearance.Id = s.AppearanceId
	metachat.Appearance.Title = s.AppearanceTitle
	metachat.Appearance.Layout = s.AppearanceLayout
	metachat.Appearance.LayoutHeaderColor = s.AppearanceLayoutHeaderColor
	metachat.Appearance.LayoutHeaderText = s.AppearanceLayoutHeaderText
	metachat.Appearance.LayoutBackGround = s.AppearanceLayoutBackGround
	metachat.Appearance.LayoutFooter = s.AppearanceLayoutFooter
	metachat.Appearance.LayoutOperatorBubble = s.AppearanceLayoutOperatorBubble
	metachat.Appearance.LayoutOperatorBubbleText = s.AppearanceLayoutOperatorBubbleText
	metachat.Appearance.LayoutVisitorBubble = s.AppearanceLayoutVisitorBubble
	metachat.Appearance.LayoutVisitorBubbleText = s.AppearanceLayoutVisitorBubbleText
	metachat.Appearance.BtnPosition = s.AppearanceBtnPosition
	metachat.Appearance.ShowAvatar = s.AppearanceShowAvatar
	metachat.Appearance.AvatarId = s.AppearanceAvatarId
	metachat.Appearance.IconButton = s.AppearanceIconButton
	metachat.Appearance.ShowPositionNum = s.AppearanceShowPositionNum
	metachat.Appearance.PlaySoundStart = s.AppearancePlaySoundStart
	metachat.Appearance.PlaySoundDecre = s.AppearancePlaySoundDecre
	metachat.Appearance.OfferTalkBot = s.AppearanceOfferTalkBot
	//metachat.Appearance.Logo = s.AppearanceLogo

	metachat.Deployment = &Deployment{}
	metachat.Deployment.Id = s.DeploymentId
	metachat.Deployment.DeployHtmlCode = s.DeploymentDeployHtmlCode
	metachat.Deployment.MailToSendHtmlCode = s.DeploymentMailToSendHtmlCode
	metachat.Deployment.DomainAllowedDeployment = s.DeploymentDomainAllowedDeployment

	metachat.Prechat = &Prechat{}
	metachat.Prechat.Id = s.PrechatId
	metachat.Prechat.Enabled = s.PrechatEnabled
	metachat.Prechat.Name = s.PrechatName
	metachat.Prechat.NameRequired = s.PrechatNameRequired
	metachat.Prechat.Lastname = s.PrechatLastname
	metachat.Prechat.LastnameRequired = s.PrechatLastnameRequired
	metachat.Prechat.Email = s.PrechatEmail
	metachat.Prechat.EmailRequired = s.PrechatEmailRequired
	metachat.Prechat.Question = s.PrechatQuestion
	metachat.Prechat.QuestionRequired = s.PrechatQuestionRequired
	metachat.Prechat.Custom1 = s.PrechatCustom1
	metachat.Prechat.Custom1Required = s.PrechatCustom1Required
	metachat.Prechat.Custom2 = s.PrechatCustom2
	metachat.Prechat.Custom2Required = s.PrechatCustom2Required
	metachat.Prechat.Custom3 = s.PrechatCustom3
	metachat.Prechat.Custom3Required = s.PrechatCustom3Required

	metachat.Postchat = &Postchat{}
	metachat.Postchat.Id = s.PostchatId
	metachat.Postchat.Enabled = s.PostchatEnabled
	metachat.Postchat.SendSurvey = s.PostchatSendSurvey
	metachat.Postchat.Ask = s.PostchatAsk
	metachat.Postchat.ShowSolved = s.PostchatShowSolved
	metachat.Postchat.SendEndMsg = s.PostchatSendEndMsg
}

// Map Metachat from single row
func (metachat *Metachat) MapMetachat(row *sql.Row) error {

	var s MetachatRow
	err := row.Scan(&s.MetachatId, &s.MetachatCustomerId, &s.MetachatSiteId, &s.MetachatThemeId, &s.MetachatCreationDate, &s.MetachatIsEnabled, &s.MetachatIsDeleted,
		&s.AppearanceId, &s.AppearanceTitle, &s.AppearanceLayout, &s.AppearanceLayoutHeaderColor, &s.AppearanceLayoutHeaderText, &s.AppearanceLayoutBackGround, &s.AppearanceLayoutFooter, &s.AppearanceLayoutOperatorBubble,
		&s.AppearanceLayoutOperatorBubbleText, &s.AppearanceLayoutVisitorBubble, &s.AppearanceLayoutVisitorBubbleText, &s.AppearanceBtnPosition, &s.AppearanceShowAvatar, &s.AppearanceAvatarId,
		&s.AppearanceIconButton, &s.AppearanceShowPositionNum, &s.AppearancePlaySoundStart, &s.AppearancePlaySoundDecre, &s.AppearanceOfferTalkBot,
		&s.DeploymentId, &s.DeploymentDeployHtmlCode, &s.DeploymentMailToSendHtmlCode, &s.DeploymentDomainAllowedDeployment,
		&s.PrechatId, &s.PrechatEnabled, &s.PrechatName, &s.PrechatNameRequired, &s.PrechatLastname, &s.PrechatLastnameRequired, &s.PrechatEmail,
		&s.PrechatEmailRequired, &s.PrechatQuestion, &s.PrechatQuestionRequired, &s.PrechatCustom1, &s.PrechatCustom1Required,
		&s.PrechatCustom2, &s.PrechatCustom2Required, &s.PrechatCustom3, &s.PrechatCustom3Required,
		&s.PostchatId, &s.PostchatEnabled, &s.PostchatSendSurvey, &s.PostchatAsk, &s.PostchatShowSolved, &s.PostchatSendEndMsg)

	if err != nil {
		return err
	}

	metachat.loadFromMetachat(s)
	fmt.Println("Metachat:", metachat)

	return nil
}

// Map Metachat from rows
func (metachats *Metachats) MapMetachats(rows *sql.Rows) error {

	for rows.Next() {

		var metachat Metachat
		var s MetachatRow

		err := rows.Scan(&s.MetachatId, &s.MetachatCustomerId, &s.MetachatSiteId, &s.MetachatThemeId, &s.MetachatCreationDate, &s.MetachatIsEnabled, &s.MetachatIsDeleted,
			&s.AppearanceId, &s.AppearanceTitle, &s.AppearanceLayout, &s.AppearanceLayoutHeaderColor, &s.AppearanceLayoutHeaderText, &s.AppearanceLayoutBackGround, &s.AppearanceLayoutFooter, &s.AppearanceLayoutOperatorBubble,
			&s.AppearanceLayoutOperatorBubbleText, &s.AppearanceLayoutVisitorBubble, &s.AppearanceLayoutVisitorBubbleText, &s.AppearanceBtnPosition, &s.AppearanceShowAvatar, &s.AppearanceAvatarId,
			&s.AppearanceIconButton, &s.AppearanceShowPositionNum, &s.AppearancePlaySoundStart, &s.AppearancePlaySoundDecre, &s.AppearanceOfferTalkBot,
			&s.DeploymentId, &s.DeploymentDeployHtmlCode, &s.DeploymentMailToSendHtmlCode, &s.DeploymentDomainAllowedDeployment,
			&s.PrechatId, &s.PrechatEnabled, &s.PrechatName, &s.PrechatNameRequired, &s.PrechatLastname, &s.PrechatLastnameRequired, &s.PrechatEmail,
			&s.PrechatEmailRequired, &s.PrechatQuestion, &s.PrechatQuestionRequired, &s.PrechatCustom1, &s.PrechatCustom1Required,
			&s.PrechatCustom2, &s.PrechatCustom2Required, &s.PrechatCustom3, &s.PrechatCustom3Required,
			&s.PostchatId, &s.PostchatEnabled, &s.PostchatSendSurvey, &s.PostchatAsk, &s.PostchatShowSolved, &s.PostchatSendEndMsg)

		if err != nil {
			return err
		}

		metachat.loadFromMetachat(s)
		fmt.Println("Metachat:", metachat)

		*metachats = append(*metachats, &metachat)
	}

	return nil
}

func (metachat *Metachat) FromSite(site Site) {

	if metachat.SiteId != 0 {

		metachat.SiteTextDescrip = site.Name + " - " + site.Domain
	}
}

func (metachats Metachats) MapSites(sites []Site) Metachats {

	var result Metachats = make(Metachats, 0)

	for _, metachat := range metachats {

		for _, site := range sites {

			if metachat.SiteId == site.Id {

				metachat.FromSite(site)
				result = append(result, metachat)
			}
		}
	}
	return result
}
