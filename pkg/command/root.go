// Copyright © 2022 zc2638 <zc2638@qq.com>.

package command

import (
	"errors"
	"fmt"
	"sync"

	"github.com/panjf2000/ants/v2"
	"github.com/spf13/cobra"
	"github.com/zc2638/ylgy/pkg/core"
)

type Option struct {
	UID   string
	Times int
}

func NewRootCommand() *cobra.Command {
	opt := new(Option)

	cmd := &cobra.Command{
		Use:          "ylgy",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(opt.UID) == 0 {
				return errors.New("uid 必填")
			}

			times := 1
			if opt.Times > 0 {
				times = opt.Times
			}

			var wg sync.WaitGroup
			p, err := ants.NewPoolWithFunc(1000, func(i interface{}) {
				defer wg.Done()

				if err := core.SendByUID(opt.UID); err != nil {
					fmt.Printf("[%d] 请求错误: %v \n", i, err)
				} else {
					fmt.Printf("[%d] 通关！\n", i)
				}
			})
			if err != nil {
				return err
			}
			defer p.Release()

			for i := 1; i <= times; i++ {
				wg.Add(1)
				_ = p.Invoke(i)
			}
			wg.Wait()
			return nil
		},
	}

	cmd.Flags().StringVar(&opt.UID, "uid", opt.UID, "设置uid")
	cmd.Flags().IntVar(&opt.Times, "times", opt.Times, "设置次数")
	return cmd
}
