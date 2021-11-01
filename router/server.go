package router

import (
	"fmt"
	"io/ioutil"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	"github.com/ultimatepanel2000/backend/db"
)

func InitServer(r *gin.RouterGroup) {

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	r.GET("/test", func(c *gin.Context) {
		go func() {

			resp, err := cli.ContainerCreate(ctx, &container.Config{
				Image:     "debian",
				OpenStdin: true,
			}, &container.HostConfig{
				Resources: container.Resources{
					Memory:            128 * 1024 * 1024,
					MemoryReservation: 80 * 1024 * 1024,
					CPUShares:         256,
				},
			}, nil, nil, "test2")
			if err != nil {
				panic(err)
			}

			if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
				panic(err)
			}

			r, _ := cli.ContainerAttach(ctx, resp.ID, types.ContainerAttachOptions{
				Stdin:  true,
				Stdout: true,
				Stderr: true,
			})

			r.Conn.Write([]byte("ls"))

			b, _ := ioutil.ReadAll(r.Reader)

			fmt.Println(string(b))

			fmt.Println(resp.ID)
		}()
		c.String(200, "ok")
	})

	r.Use(NeedsAuth)

	r.GET("/", func(c *gin.Context) {
		c.String(200, "ok")
	})

	r.PUT("/new/:name", func(c *gin.Context) {
		s, err := db.DB.Server.CreateOne(
			db.Server.Owner.Link(
				db.User.ID.Equals(c.GetString("user_id")),
			),
			db.Server.Name.Set(c.Param("name")),
			db.Server.Status.Set(0),
		).Exec(ctx)

		if handleError(err, c, "db error") {
			return
		}

		c.JSON(200, gin.H{
			"status": "success",
			"server": s,
		})
	})
}
