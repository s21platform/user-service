package new_avatar

type DBRepo interface {
	UpdateUserAvatar(UUID, link string) error
}
