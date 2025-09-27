package ssh

// func tunneling(ctx context.Context, client *ssh.Client) error {
// 	if err := promptTunneling(config); err != nil {
// 		return err
// 	}

// 	lsnr, err := net.Listen("tcp", "127.0.0.1:"+config.LocalPort)
// 	if err != nil {
// 		return err
// 	}
// 	defer lsnr.Close()

// 	// spinner()

// 	var wg sync.WaitGroup
// 	var errCh = make(chan error, 2)
// 	go func(ctx context.Context) {
// 		fmt.Printf("Local Port(:%s) Listening\n", config.LocalPort)
// 		localConn, err := lsnr.Accept()
// 		if err != nil {
// 			errCh <- err
// 			close(errCh)
// 			return
// 		}
// 		defer localConn.Close()

// 		remoteConn, err := client.Dial("tcp", "127.0.0.1:"+config.RemotePort)
// 		if err != nil {
// 			errCh <- err
// 			close(errCh)
// 			return
// 		}
// 		defer remoteConn.Close()

// 		wg.Add(1)
// 		go func() {
// 			defer wg.Done()
// 			if _, err := CopyContext(ctx, remoteConn, localConn); err != nil {
// 				errCh <- fmt.Errorf("remote <--- local, '%w'", err)
// 			}
// 		}()

// 		wg.Add(1)
// 		go func() {
// 			defer wg.Done()
// 			if _, err := CopyContext(ctx, localConn, remoteConn); err != nil {
// 				errCh <- fmt.Errorf("local <--- remote, '%w'", err)
// 			}
// 		}()

// 		fmt.Printf("Local Port(:%s) <-------[tunneling established]-------> Remote Port(:%s)\r", config.LocalPort, config.RemotePort)
// 		wg.Wait()
// 		close(errCh)
// 	}(ctx)

// 	var firstError error
// 	var ctxErr error

// 	for {
// 		select {
// 		case <-ctx.Done():
// 			fmt.Printf("[signal] shutdown\n")
// 			ctxErr = ctx.Err()
// 			lsnr.Close()
// 		case err, ok := <-errCh:
// 			if !ok {
// 				fmt.Printf("Local Port(:%s) X-------[tunneling disconnected]-------X Remote Port(:%s)\r", config.LocalPort, config.RemotePort)
// 				if firstError != nil {
// 					return firstError
// 				}
// 				return ctxErr
// 			}
// 			if err != nil {
// 				// fmt.Printf("tunneling error: %s\n", err)
// 				if firstError == nil {
// 					firstError = err
// 				}
// 			}
// 		}
// 	}
// }

// func promptTunneling(config *Config) error {
// 	prompts := []promptui.Prompt{
// 		{
// 			Label: "Local Port",
// 		},
// 		{
// 			Label: "Remote Port",
// 		},
// 	}

// 	for i, prompt := range prompts {
// 		port, err := prompt.Run()
// 		if err != nil {
// 			return err
// 		}

// 		switch i {
// 		case 0:
// 			config.LocalPort = port
// 		case 1:
// 			config.RemotePort = port
// 		}
// 	}

// 	return nil
// }

// func CopyContext(ctx context.Context, dsc io.Writer, src io.Reader) (written int64, err error) {
// 	buf := make([]byte, 1024*32)

// 	type result struct {
// 		n   int
// 		err error
// 	}
// 	resultCh := make(chan result, 1)

// 	go func() {
// 		for {
// 			nr, readErr := src.Read(buf)
// 			select {
// 			case resultCh <- result{nr, readErr}:
// 			case <-ctx.Done():
// 				return
// 			}
// 			if readErr != nil {
// 				return
// 			}
// 		}
// 	}()

// 	for {
// 		select {
// 		case <-ctx.Done():
// 			return written, err
// 		case res := <-resultCh:
// 			if res.n > 0 {
// 				nw, writeErr := dsc.Write(buf[:res.n])
// 				if nw > 0 {
// 					written += int64(nw)
// 				}
// 				if writeErr != nil {
// 					return written, writeErr
// 				}
// 				if int(res.n) != nw {
// 					return written, io.ErrShortWrite
// 				}
// 			}
// 			if res.err != nil {
// 				if res.err == io.EOF {
// 					return written, nil
// 				}
// 				return written, res.err
// 			}
// 		}
// 	}
// }

// func spinner() {
// 	spinnerText := []string{"|", "/", "-", "\\"}
// 	for i := 0; i < 30; i++ {
// 		fmt.Printf("%s\r", spinnerText[i%len(spinnerText)])
// 		time.Sleep(100 * time.Millisecond)
// 	}
// }
