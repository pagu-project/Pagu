package repository

type Database interface {
	IUser
	IVoucher
	IFaucet
	IZealy
}
