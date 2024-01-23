package model

import "time"

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
