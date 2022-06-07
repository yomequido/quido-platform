package models

import "database/sql"

type Message struct {
	MessageId          int            `db:"message_id"`
	ExternalMessagesId sql.NullString `db:"external_message_id,omitempty"`
	SentByUser         bool           `db:"sent_by_user"`
	Channel            string         `db:"channel"`
	ChannelUserId      sql.NullString `db:"channel_user_id,omitempty"`
	FkPatient          int            `db:"fk_patient"`
	FkEmployee         sql.NullInt64  `db:"fk_employee,omitempty"`
	Title              sql.NullString `db:"title,omitempty"`
	Message            sql.NullString `db:"message,omitempty"`
	UrlReference       sql.NullString `db:"url_reference,omitempty"`
	SentDate           sql.NullTime   `db:"sentDate,omitempty"`
	ReadDate           sql.NullTime   `db:"readDate,omitempty"`
}
