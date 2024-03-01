import React, { useEffect, useState } from "react";
import { useNavigate, useOutletContext, useParams } from "react-router-dom"
import Input from "./form/Input";
import Select from "./form/Select"
import TextArea from "./form/TextArea"
import Checkbox from "./form/Checkbox"
import Swal from 'sweetalert2'

const EditMovie = () => {
    const navigate = useNavigate()
    const { jwtToken } = useOutletContext()

    const [err, setErr] = useState(null)
    const [errs, setErrs] = useState([])

    const mpaaOptions = [
        {id: "G", value: "G"},
        {id: "PG", value: "PG"},
        {id: "PG-13", value: "PG-13"},
        {id: "R", value: "R"},
        {id: "NC-17", value: "NC-17"},
        {id: "18A", value: "18A"},
    ]

    const hasError = (key) => {
        return errs.indexOf(key) !== -1
    }

    const [movie, setMovie] = useState({
        id: 0,
        title: "",
        release_date: "",
        runtime: "",
        mpaa_rating: "",
        description: "",
        genres: [],
        genres_array: [Array(13).fill(false)],
    })

    // get id from URL
    let {id} = useParams()
    if (id === undefined) {
        id = 0
    }

    useEffect(() => {
        if (jwtToken === "") {
            navigate("/login")
            return
        }

        if (id === 0) {
            // add a movie
            setMovie({
                id: 0,
                title: "",
                release_date: "",
                runtime: "",
                mpaa_rating: "",
                description: "",
                genres: [],
                genres_array: [Array(13).fill(false)],
            })

            const headers = new Headers()
            headers.append("Content-Type", "application/json")

            const requestOptions = {
                method: "GET",
                headers: headers,
            }

            fetch(`${process.env.REACT_APP_BACKEND}/genres`, requestOptions)
                .then(response => response.json())
                .then(data => {
                    const checks = []

                    data.forEach(g => {
                        checks.push({id: g.id, checked: false, genre: g.genre})
                    })

                    setMovie(m => ({
                        ...m,
                        genres: checks,
                        genres_array: [],
                    }))
                })
                .catch(err => console.log(err))
        } else {
            // edit an existing movie
            const headers = new Headers()
            headers.append("Content-Type", "application/json")
            headers.append("Authorization", "Bearer " + jwtToken)

            const requestOptions = {
                method: "GET",
                headers: headers,
            }

            fetch(`${process.env.REACT_APP_BACKEND}/admin/movies/${id}`, requestOptions)
                .then(response => {
                    if (response.status !== 200) {
                        setErr("Invalid response code: " + response.status)
                    }
                    return response.json()
                .then(data => {
                    // fix release date
                    data.movie.release_date = new Date(data.movie.release_date).toISOString().split('T')[0]

                    const checks = []

                    data.genres.forEach(g => {
                        if (data.movie.genres_array.indexOf(g.id) !== -1) {
                            checks.push({id: g.id, checked: true, genre: g.genre})
                        } else {
                            checks.push({id: g.id, checked: false, genre: g.genre})
                        }
                    })

                    // set state
                    setMovie({
                        ...data.movie,
                        genres: checks,
                    })
                })
                .catch(err => console.log(err))
            })
        }
    }, [id, jwtToken])

    const handleSubmit = (event) => {
        event.preventDefault()

        let errors = []
        let required = [
            { field: movie.title, name: "title" },
            { field: movie.release_date, name: "release_date" },
            { field: movie.runtime, name: "runtime" },
            { field: movie.description, name: "description" },
            { field: movie.mpaa_rating, name: "mpaa_rating" },
        ]

        required.forEach(obj => {
            if (obj.field === "") {
                errors.push(obj.name)
            }
        })

        if (movie.genres_array.length === 0) {
            Swal.fire({
                title: 'Error!',
                text: 'You must choose genre',
                icon: 'error',
                confirmButtonText: 'OK',
            })
            errors.push("genres")
        }

        setErrs(errors)

        if (errors.length > 0) {
            return false
        }

        // passed validation, then save changes
        const headers = new Headers()
        headers.append("Content-Type", "application/json")
        headers.append("Authorization", "Bearer " + jwtToken)

        // assume we r adding new movie
        let method = "POST"

        if (movie.id > 0) {
            method = "PATCH"
        }

        const requestBody = movie
        // convert values in JSON 4 release_date (to date)
        // runtime to int

        requestBody.release_date = new Date(movie.release_date)
        requestBody.runtime = parseInt(movie.runtime, 10)

        let requestOptions = {
            body: JSON.stringify(requestBody),
            method: method,
            headers: headers,
            credentials: "include",
        }

        fetch(`${process.env.REACT_APP_BACKEND}/admin/movies/${movie.id}`, requestOptions)
            .then(response => response.json())
            .then(data => {
                if (data.error) {
                    console.log(data.message)
                } else {
                    navigate("/manage-catalog")
                }
            })
            .catch(err => console.log(err))
    }

    const handleChange = () => (event) => {
        let value = event.target.value
        let name = event.target.name
        setMovie({
            ...movie,
            [name]: value,
        })
    }

    const handleCheck = (event, position) => {
        console.log("handleCheck called")
        console.log("handleCheck:", event.target.value)
        console.log("checked is", event.target.checked)
        console.log("position is", position)

        let tmpArr = movie.genres
        tmpArr[position].checked = !tmpArr[position].checked  //reverse the value

        let valueInt = parseInt(event.target.value, 10)

        let tmpIDs = movie.genres_array
        if (!event.target.checked) {
            tmpIDs.splice(tmpIDs.indexOf(valueInt), 1)
        } else {
            tmpIDs.push(valueInt)
        }

        setMovie({
            ...movie,
            genres: tmpArr,
            genres_array: tmpIDs,
        })
    }

    const confirmDelete = () => {
        Swal.fire({
            title: "Delete movie?",
            text: "You can't undo this!",
            icon: "warning",
            showCancelButton: true,
            confirmButtonColor: "#3085d6",
            cancelButtonColor: "#d33",
            confirmButtonText: "Yes, delete it!"
          }).then((result) => {
            if (result.isConfirmed) {
                let headers = new Headers()
                headers.append("Authorization", "Bearer " + jwtToken)

                const requestOptions = {
                    method: "DELETE",
                    headers: headers,
                }
                fetch(`${process.env.REACT_APP_BACKEND}/admin/movies/${movie.id}`, requestOptions)
                    .then(response => response.json())
                    .then(data => {
                        if (data.error) {
                            console.log(data.message)
                        } else {
                            navigate("/manage-catalog")
                        }
                    })
                    .catch(err => console.log(err))
            }
          });
    }
    
    if (err !== null) {
        return <div>Error: {err.message}</div>
    } else {
        return (
        <div>
            <h2>Add/Edit Movie</h2>
            <hr />
            {/* <pre>{JSON.stringify(movie, null, 3)}</pre> */}

            <form onSubmit={handleSubmit}>
                <input type="hidden" name="id" value={movie.id} id="id"></input>
                <Input
                    title={"Title"}
                    className={"form-control"}
                    type={"text"}
                    name={"title"}
                    value={movie.title}
                    onChange={handleChange("title")}
                    errorDiv={hasError("title") ? "text-danger" : "d-none"}
                    errorMsg={"Please enter a title"}
                />

                <Input
                    title={"Release Date"}
                    className={"form-control"}
                    type={"date"}
                    name={"release_date"}
                    value={movie.release_date}
                    onChange={handleChange("release_date")}
                    errorDiv={hasError("release_date") ? "text-danger" : "d-none"}
                    errorMsg={"Please enter a release date"}
                />

                <Input
                    title={"Runtime"}
                    className={"form-control"}
                    type={"text"}
                    name={"runtime"}
                    value={movie.runtime}
                    onChange={handleChange("runtime")}
                    errorDiv={hasError("runtime") ? "text-danger" : "d-none"}
                    errorMsg={"Please enter a runtime"}
                />

                <Select
                    title={"MPAA Rating"}
                    name={"mpaa_rating"}
                    options={mpaaOptions}
                    value={movie.mpaa_rating}
                    onChange={handleChange("mpaa_rating")}
                    placeHolder={"Choose..."}
                    errorMsg={"Please choose"}
                    errorDiv={hasError("mpaa_rating") ? "text-danger" : "d-none"}
                />

                <TextArea
                    title="Description"
                    name={"description"}
                    value={movie.description}
                    rows={"3"}
                    onChange={handleChange("description")}
                    errorDiv={hasError("description") ? "text-danger" : "d-none"}
                    errorMsg={"Please enter a description"}
                />

                <hr />
                <h3>Genres</h3>

                {movie.genres && movie.genres.length > 1 && 
                    <>
                        {Array.from(movie.genres).map((g, index) => (
                            <Checkbox 
                                title={g.genre}
                                name={"genre"}
                                key={index}
                                id={"genre-"+index}
                                onChange={event => handleCheck(event, index)}
                                value={g.id}
                                checked={movie.genres[index].checked}
                            />
                        ))}
                    </>
                }

                <hr />

                <button className="btn btn-primary">Save</button>
                
                {movie.id > 0 &&
                    <a href="#" className="btn btn-danger ms-2" onClick={confirmDelete}>Delete Movie</a>
                }
            </form>
        </div>
    )}
}

export default EditMovie