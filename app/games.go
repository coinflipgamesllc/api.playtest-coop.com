package app

// func (s *Server) handleListGames() gin.HandlerFunc {
// 	type request struct {
// 		Title       string `form:"title"`
// 		Status      string `form:"status"`
// 		Designer    string `form:"designer"`
// 		PlayerCount int    `form:"player_count"`
// 		Age         int    `form:"age"`
// 		Playtime    int    `form:"playtime"`
// 		Limit       int    `form:"limit"`
// 		Offset      int    `form:"offset"`
// 		Sort        string `form:"sort"`
// 	}

// 	type response struct {
// 		Games  []domain.Game `json:"games"`
// 		Total  int           `json:"total"`
// 		Limit  int           `json:"limit"`
// 		Offset int           `json:"offset"`
// 	}

// 	return func(c *gin.Context) {
// 		// Validate request
// 		var req request
// 		if err := c.ShouldBind(&req); err != nil {
// 			c.AbortWithStatusJSON(400, serverError(err))
// 			return
// 		}

// 		if req.Limit == 0 {
// 			req.Limit = 10
// 		}

// 		// Fetch games
// 		games, total, err := s.gameRepository.ListGames(
// 			req.Title,
// 			req.Status,
// 			req.Designer,
// 			req.PlayerCount,
// 			req.Age,
// 			req.Playtime,
// 			req.Limit,
// 			req.Offset,
// 			req.Sort,
// 		)

// 		if err != nil {
// 			c.AbortWithStatusJSON(500, serverError(err))
// 			return
// 		}

// 		c.JSON(200, response{Games: games, Total: total, Limit: req.Limit, Offset: req.Offset})
// 	}
// }

// func (s *Server) handleCreateGame() gin.HandlerFunc {
// 	type stats struct {
// 		MinPlayers        int `json:"min_players" binding:"min=0,ltefield=MaxPlayers"`
// 		MaxPlayers        int `json:"max_players" binding:"min=0,gtefield=MinPlayers"`
// 		MinAge            int `json:"min_age" binding:"min=0,max=99"`
// 		EstimatedPlaytime int `json:"estimated_playtime" binding:"min=0,max=9999"`
// 	}

// 	type request struct {
// 		Title     string `json:"title" binding:"required"`
// 		Overview  string `json:"overview"`
// 		Designers []uint `json:"designers"`
// 		Stats     *stats `json:"stats" binding:"omitempty,dive"`
// 	}

// 	type response struct {
// 		Game *domain.Game `json:"game"`
// 	}

// 	return func(c *gin.Context) {
// 		// Validate request
// 		var req request
// 		if err := c.ShouldBind(&req); err != nil {
// 			c.AbortWithStatusJSON(400, serverError(err))
// 			return
// 		}

// 		// Create our new game
// 		currentUser := s.user(c)
// 		game := domain.NewGame(req.Title, *currentUser)

// 		// If the request included optional information, add it now
// 		if req.Overview != "" {
// 			game.UpdateOverview(req.Overview)
// 		}

// 		if len(req.Designers) > 1 { // Index 0 is always the current user, which is included already
// 			for _, designerID := range req.Designers {
// 				if designerID == currentUser.ID {
// 					continue
// 				}

// 				designer, err := s.userRepository.UserOfID(designerID)
// 				if err != nil {
// 					c.AbortWithStatusJSON(500, serverError(err))
// 					return
// 				}

// 				game.AddDesigner(designer)
// 			}
// 		}

// 		if req.Stats != nil {
// 			game.UpdateStats(req.Stats.MinPlayers, req.Stats.MaxPlayers, req.Stats.MinAge, req.Stats.EstimatedPlaytime)
// 		}

// 		// And save
// 		err := s.gameRepository.Save(game)
// 		if err != nil {
// 			c.AbortWithStatusJSON(500, serverError(err))
// 			return
// 		}

// 		c.JSON(200, response{Game: game})
// 	}
// }

// func (s *Server) handleGetGame() gin.HandlerFunc {
// 	type response struct {
// 		Game *domain.Game `json:"game"`
// 	}

// 	return func(c *gin.Context) {
// 		// Validate request
// 		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
// 		if err != nil {
// 			c.AbortWithStatusJSON(500, serverError(err))
// 			return
// 		}

// 		game, err := s.gameRepository.GameOfID(uint(id))
// 		if err != nil {
// 			c.AbortWithStatusJSON(500, serverError(err))
// 			return
// 		}

// 		if game == nil {
// 			c.AbortWithStatusJSON(404, serverError(errors.New("not found")))
// 			return
// 		}

// 		c.JSON(200, response{Game: game})
// 	}
// }

// func (s *Server) handleUpdateGame() gin.HandlerFunc {
// 	type stats struct {
// 		MinPlayers        int `json:"min_players" binding:"min=0,ltefield=MaxPlayers"`
// 		MaxPlayers        int `json:"max_players" binding:"min=0,gtefield=MinPlayers"`
// 		MinAge            int `json:"min_age" binding:"min=0,max=99"`
// 		EstimatedPlaytime int `json:"estimated_playtime" binding:"min=0,max=9999"`
// 	}

// 	type request struct {
// 		Title     string `json:"title"`
// 		Overview  string `json:"overview"`
// 		Status    string `json:"status"`
// 		Designers []uint `json:"designers"`
// 		Stats     *stats `json:"stats" binding:"omitempty,dive"`
// 	}

// 	type response struct {
// 		Game *domain.Game `json:"game"`
// 	}

// 	return func(c *gin.Context) {
// 		// Pull game by ID
// 		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
// 		if err != nil {
// 			c.AbortWithStatusJSON(500, serverError(err))
// 			return
// 		}

// 		game, err := s.gameRepository.GameOfID(uint(id))
// 		if err != nil {
// 			c.AbortWithStatusJSON(500, serverError(err))
// 			return
// 		}

// 		if game == nil {
// 			c.AbortWithStatusJSON(404, serverError(errors.New("not found")))
// 			return
// 		}

// 		// Ensure that our current user is allowed to edit the game
// 		currentUser := s.user(c)
// 		if !game.MayBeUpdatedBy(currentUser) {
// 			c.AbortWithStatusJSON(401, serverError(errors.New("you may not edit this game")))
// 			return
// 		}

// 		// Validate the request itself
// 		var req request
// 		if err := c.ShouldBind(&req); err != nil {
// 			c.AbortWithStatusJSON(400, serverError(err))
// 			return
// 		}

// 		// Update game
// 		if req.Title != "" {
// 			game.Rename(req.Title)
// 		}

// 		if req.Overview != "" {
// 			game.UpdateOverview(req.Overview)
// 		}

// 		if req.Status != "" {
// 			err := game.UpdateStatus(req.Status)
// 			if err != nil {
// 				c.AbortWithStatusJSON(400, serverError(err))
// 				return
// 			}
// 		}

// 		if len(req.Designers) > 0 {
// 			designers := []domain.User{}
// 			for _, designerID := range req.Designers {
// 				designer, err := s.userRepository.UserOfID(designerID)
// 				if err != nil {
// 					c.AbortWithStatusJSON(500, serverError(err))
// 					return
// 				}

// 				designers = append(designers, *designer)
// 			}

// 			game.ReplaceDesigners(designers)
// 		}

// 		if req.Stats != nil {
// 			game.UpdateStats(req.Stats.MinPlayers, req.Stats.MaxPlayers, req.Stats.MinAge, req.Stats.EstimatedPlaytime)
// 		}

// 		// And save
// 		err = s.gameRepository.Save(game)
// 		if err != nil {
// 			c.AbortWithStatusJSON(500, serverError(err))
// 			return
// 		}

// 		c.JSON(200, response{Game: game})
// 	}
// }
