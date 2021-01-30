package database

func GetAllChannel() (channels []Channel, err error) {
	err = db.Find(&channels).Error
	return
}

func SaveChannel(channel Channel) error {
	return db.Save(&channel).Error
}

func DeleteChannel(id uint) error {
	return db.Delete(Channel{}, "id = ?", id).Error
}

func GetChannel(channelNumber uint) (channel Channel, err error) {
	err = db.Where("id = ?", channelNumber).First(&channel, channelNumber).Error
	return
}
