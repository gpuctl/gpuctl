package gpustats

type GPUDataSource interface {
	GetGPUStatus() (GPUStatusPacket, error)
}
