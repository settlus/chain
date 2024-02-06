package types

// ValidateBasic is used for validating the packet
func (p ForwardingAccountPacket) ValidateBasic() error {
	addr := p.Registrant
	if addr == "" {
		return ErrInvalidPacket
	}

	// TODO: test valid bech32 address

	return nil
}

// GetBytes is a helper for serialising
func (p ForwardingAccountPacket) GetBytes() ([]byte, error) {
	var modulePacket ForwardingAccountPacketData
	modulePacket.Packet = &ForwardingAccountPacketData_ForwardingAccountPacket{ForwardingAccountPacket: &p}
	return modulePacket.Marshal()
}
