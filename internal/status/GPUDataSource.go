package status

type GPUDataSource interface {
	GetGPUStatus() (GPUStatusPacket, error)
}
