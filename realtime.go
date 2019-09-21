package bitflyer

type Board struct{}

func (r *RealtimeAPIClient) GetBoard() (*Board, error) {
	return &Board{}, nil
}

func (r *RealtimeAPIClient) Subscribe() error {
	if err := r.rpc.Open(); err != nil {
		return err
	}

	r.rpc.Recv()
	return nil
}

func (r *RealtimeAPIClient) Close() error {
	return r.rpc.Close()
}
