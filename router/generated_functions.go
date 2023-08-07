package router

// postPostGenerated godoc
// @Summary   Creates entity
// @Tags   post crud
// @Accept    json
// @Produce   json
// @Param	  request body post.CreateDto true "Request data"
// @Response  201 {object}        post.Model
// @Router    /campaign/crud/:campaign/post [post]
func postPostGenerated() {}

// postUpdateGenerated godoc
// @Summary   Updates entity
// @Tags   post crud
// @Accept    json
// @Produce   json
// @Param	   id path string true "id of item"
// @Param	  request body post.UpdateDto true "Request data"
// @Response  200 {object}        post.Model
// @Router    /campaign/crud/:campaign/post/:id [put]
func postUpdateGenerated() {}

// postDeleteGenerated godoc
// @Summary   Deletes entity
// @Tags   post crud
// @Accept    json
// @Produce   json
// @Param	   id path string true "id of item"
// @Router    /campaign/crud/:campaign/post/:id [delete]
func postDeleteGenerated() {}

// postGetGenerated godoc
// @Summary   Returns entity
// @Tags   post crud
// @Accept    json
// @Produce   json
// @Param	   id path string true "id of item"
// @Response  200 {object}        post.Model
// @Router    /campaign/crud/:campaign/post/:id [get]
func postGetGenerated() {}

// postListGenerated godoc
// @Summary   Lists entities
// @Tags   post crud
// @Accept    json
// @Produce   json
// @Param     campaign            path string true  "campaign uuid"
// @Param     offset              query string false  "offset"
// @Param     limit               query string false  "limit"
// @Param     orderBy             query string false  "orderBy"
// @Param     orderDir            query string false  "orderDir"
// @Response  200 {object}        pagination.Output{rows=[]post.Model}
// @Router    /campaign/crud/:campaign/post [get]
func postListGenerated() {}

// authPostGenerated godoc
// @Summary   Creates entity
// @Tags   auth crud
// @Accept    json
// @Produce   json
// @Param	  request body auth.CreateDto true "Request data"
// @Response  201 {object}        auth.Model
// @Router    /campaign/crud/:campaign/auth [post]
func authPostGenerated() {}

// authUpdateGenerated godoc
// @Summary   Updates entity
// @Tags   auth crud
// @Accept    json
// @Produce   json
// @Param	   id path string true "id of item"
// @Param	  request body auth.UpdateDto true "Request data"
// @Response  200 {object}        auth.Model
// @Router    /campaign/crud/:campaign/auth/:id [put]
func authUpdateGenerated() {}

// authDeleteGenerated godoc
// @Summary   Deletes entity
// @Tags   auth crud
// @Accept    json
// @Produce   json
// @Param	   id path string true "id of item"
// @Router    /campaign/crud/:campaign/auth/:id [delete]
func authDeleteGenerated() {}

// authGetGenerated godoc
// @Summary   Returns entity
// @Tags   auth crud
// @Accept    json
// @Produce   json
// @Param	   id path string true "id of item"
// @Response  200 {object}        auth.Model
// @Router    /campaign/crud/:campaign/auth/:id [get]
func authGetGenerated() {}

// authListGenerated godoc
// @Summary   Lists entities
// @Tags   auth crud
// @Accept    json
// @Produce   json
// @Param     campaign            path string true  "campaign uuid"
// @Param     offset              query string false  "offset"
// @Param     limit               query string false  "limit"
// @Param     orderBy             query string false  "orderBy"
// @Param     orderDir            query string false  "orderDir"
// @Response  200 {object}        pagination.Output{rows=[]auth.Model}
// @Router    /campaign/crud/:campaign/auth [get]
func authListGenerated() {}
