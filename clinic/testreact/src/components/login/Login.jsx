import React, { useState } from "react";
import { Link } from "react-router-dom";

import './Login.css';
import axios from "axios";


const Login = () => {
    const [formData, setFormData] = useState({ username: '', password: '' });
    const [error, setError] = useState('');

    const handleFormSubmit = () => {
        axios.post(`http://localhost:12345/SignIN`, formData)
            .then(response => {
                console.log(response.data);
                if (response.status === 200) {
                    const userType = response.data.typeee;

                    // Store token and user type in local storage
                    localStorage.setItem('token', response.data.token);
                    localStorage.setItem('userType', userType);
                    localStorage.setItem('username', formData.username);

                    setError("Signed In Successfully");

                    window.open(`/${userType}`, '_blank');

                    window.location.reload();
                }
            })
            .catch(error => {
                if (error.response) {
                    if (error.response.status === 400) {
                        setError('Invalid credentials,please enter again Password or Username .');
                    }else if (error.response.status === 404) {
                        setError('please enter the missing info.');
                    }
                    else if (error.response.status === 411) {
                        setError('the length of email is greater than 50 or password is greater than 20.');
                    }
                    else if (error.response.status === 406) {
                        setError('email must be formatted as (a-zA-Z0-9._%+-)@(a-zA-Z0-9._%+-).2characterORmore   and password must be at least 8 characters with atleast 1 digit and 1 Capital letter');
                    }
                    else if (error.response.status === 423) {
                        setError('this account dont have an type,so please create another account to determine patient/doctor');
                    }

                } else {
                    setError(' error. Please try again.');
                }
            });
    }
    return (
        <div className="container ">
            <div className="header">
                <div className="text">Log In</div>
                <div className="underline"></div>
            </div>
            {error && <div className="error-message header "  style={{ color: 'red' }}>{error}</div>}
            <form>
                <div className="inputs">

                    <div className="input">
                        <label htmlFor="username" className="info">Email</label>
                        <input type="text" placeholder="your email" value={formData.username} onChange={e => setFormData({ ...formData, username: e.target.value })}/>
                    </div>
                    <div className="input">
                        <label htmlFor="password" className="info">Password</label>
                        <input type="text" placeholder="your password" value={formData.password} onChange={e => setFormData({ ...formData, password: e.target.value })} />
                    </div>
                    <div className="submit-container">
                        <button type="button" onClick={handleFormSubmit} className="Login">Login</button>
                        <button><Link to="/Signup" type="submit" className="Signup">Sign UP</Link></button>
                    </div>
                </div>
            </form>
        </div>
    )
}
export default Login