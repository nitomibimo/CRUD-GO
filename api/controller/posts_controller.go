package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/nitomibimo/CRUD-GO/api/auth"
	"github.com/nitomibimo/CRUD-GO/api/models"
	"github.com/nitomibimo/CRUD-GO/api/response"
	"github.com/nitomibimo/CRUD-GO/api/utils/formaterror"
)

// CreatePost - Function for create article
func (s *Server) CreatePost(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	post := models.Post{}
	err = json.Unmarshal(body, &post)
	if err != nil {
		response.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	post.Prepare()
	err = post.Validate()
	if err != nil {
		response.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if uid != post.AuthorID {
		response.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	postCreated, err := post.SavePost(s.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		response.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Lacation", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, postCreated.ID))
	response.JSON(w, http.StatusCreated, postCreated)
}

//GetPosts - function for read all article
func (s *Server) GetPosts(w http.ResponseWriter, r *http.Request) {

	post := models.Post{}

	posts, err := post.FindAllPosts(s.DB)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	response.JSON(w, http.StatusOK, posts)
}

// GetPost - function for sorting article by id
func (s *Server) GetPost(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}
	post := models.Post{}

	postReceived, err := post.FindPostByID(s.DB, pid)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	response.JSON(w, http.StatusOK, postReceived)
}

// UpdatePost - function for after edit aticle
func (s *Server) UpdatePost(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// Check if the post id is valid
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	//CHeck if the auth token is valid and  get the user id from it
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// Check if the post exist
	post := models.Post{}
	err = s.DB.Debug().Model(models.Post{}).Where("id = ?", pid).Take(&post).Error
	if err != nil {
		response.ERROR(w, http.StatusNotFound, errors.New("Post not found"))
		return
	}

	// If a user attempt to update a post not belonging to him
	if uid != post.AuthorID {
		response.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	// Read the data posted
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Start processing the request data
	postUpdate := models.Post{}
	err = json.Unmarshal(body, &postUpdate)
	if err != nil {
		response.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	//Also check if the request user id is equal to the one gotten from token
	if uid != postUpdate.AuthorID {
		response.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	postUpdate.Prepare()
	err = postUpdate.Validate()
	if err != nil {
		response.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	postUpdate.ID = post.ID //this is important to tell the model the post id to update, the other update field are set above

	postUpdated, err := postUpdate.UpdateAPost(s.DB)

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		response.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	response.JSON(w, http.StatusOK, postUpdated)
}

// DeletePost - function for deleting article
func (s *Server) DeletePost(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// Is a valid post id given to us?
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Is this user authenticated?
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// Check if the post exist
	post := models.Post{}
	err = s.DB.Debug().Model(models.Post{}).Where("id = ?", pid).Take(&post).Error
	if err != nil {
		response.ERROR(w, http.StatusNotFound, errors.New("Unauthorized"))
		return
	}

	// Is the authenticated user, the owner of this post?
	if uid != post.AuthorID {
		response.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	_, err = post.DeleteAPost(s.DB, pid, uid)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", pid))
	response.JSON(w, http.StatusNoContent, "")
}
