import React, { useEffect, useState } from "react";
import { useNavigate, useOutletContext, useParams } from "react-router-dom"
import Input from "./form/Input";

const EditMovie = () => {
    const navigate = useNavigate()
    const { jwtToken } = useOutletContext()

    const [error, setError] = useState(null)
    const [errors, setErrors] = useState([])

    const hasError = (key) => {
        return errors.indexOf(key) !== -1
    }

    const [movie, setMovie] = useState({
        id: 0,
        title: "",
        release_date: "",
        runtime: "",
        mpaa_rating: "",
        description: "",
    })

    // get id from URL
    let {id} = useParams()

    useEffect(() => {
        if (jwtToken === "") {
            navigate("/login")
            return
        }
    }, [jwtToken])

    const handleSubmit = (event) => {
        event.preventDefault()
    }

    const handleChange = () => (event) => {
        let value = event.target.value
        let name = event.target.name
        setMovie({
            ...movie,
            [name]: value,
        })
    }
    
    return (
        <div>
            <h2>Add/Edit Movie</h2>
            <hr />
            <pre>{JSON.stringify(movie, null, 3)}</pre>

            <form onSubmit={handleSubmit}>
                <input type="hidden" name="id" value={movie.id} id="id"></input>
                <Input
                    title={"Title"}
                    className={"form-control"}
                    type={"text"}
                    name={"title"}
                    value={movie.title}
                    onChange={handleChange("title")}
                    errorDiv={hasError("title") ? "tex-danger" : "d-none"}
                    errorMsg={"Please enter a title"}
                />
            </form>
        </div>
    )
}

export default EditMovie