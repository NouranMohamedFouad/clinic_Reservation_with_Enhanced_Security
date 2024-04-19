import React, { useState } from "react";
import { Link } from "react-router-dom";
import './Signup.css';
import axios from 'axios';

const Signup = () => {
    const [formData, setFormData] = useState({ username: '', password: '', usertype: '', code: '' }); // Set the default user type to 'patient'
    const [error, setError]=useState('');
    const [sanitizedFormData, setSanitizedFormData] = useState({
        username: '',
        password: '',
        usertype: '',
        code: ''
    });
    const handleFormSubmit = () => {

        axios.post(`http://localhost:12345/SignUP`, sanitizedFormData)
            .then(response => {
                console.log(response.data);
                if (response.status === 200) {
                    setError('successfully signed up,please go login.');
                    window.open('/Login', '_blank');
                }
            })
            .catch(error=>{
                if (error.response.status === 400) {
                    setError('Username is already exists.');
                }
                else if (error.response.status === 404) {
                    setError('please enter the missing info.');
                }
                else if (error.response.status === 411) {
                    setError('the length of email is greater than 50 or password is greater than 20.');
                }
                else if (error.response.status === 406) {
                    setError('email must be formatted as (a-zA-Z0-9._%+-)@(a-zA-Z0-9._%+-).2characterORmore   and password must be at least 8 characters with atleast 1 digit and 1 Capital letter');
                }
                else if (error.response.status === 423) {
                    setError('Doctor , Verification Code is not correct');
                }
                else if (error.response.status === 412) {
                    setError('Verification Code must contains letters and digits only');
                }else if (error.response.status === 300) {
                    setError('Please Choose whether you are Doctor or Patient');
                }
                else {
                    setError(' error. Please try again.');
                }
            });
    }
    ///////////////////////////////////////////////////////////////////////////////////////
    const sanitizeInput = (value) => {
        return value.replace(/&/g, '&amp;')
            .replace(/</g, '&lt;')
            .replace(/>/g, '&gt;')
            .replace(/"/g, '&quot;')
            .replace(/'/g, '&#39;');
    };

    ///////////////////////////////////////////////////////////////////////////////////////


    const handleInputChange = (e, field) => {
        const { value } = e.target;
        setFormData({ ...formData, [field]: value });
        setSanitizedFormData({ ...sanitizedFormData, [field]: sanitizeInput(value) });
    };



    return(
        <div className="container ">
            <div className="header">
                <div className="text">Sign Up</div>
                <div className="underline"></div>
            </div>
            {error && <div className="error-message header "  style={{ color: 'red' }}>{error}</div>}
            <form>
                <div className="inputs">

                    <div className="input">
                        <label htmlFor="username" className="info">Email</label>
                        <input type="text" placeholder="your username" value={formData.username} onChange={e => handleInputChange(e, 'username')} />
                    </div>
                    <div className="input">
                        <label htmlFor="password" className="info">Password</label>
                        <input type="text" placeholder="your password" value={formData.password}  onChange={e => handleInputChange(e, 'password')}  />
                    </div>

                    <div className="input">
                        <label htmlFor="usertype" className="info" >User Type</label>

                        <input type="radio" className="radio"
                               id="patient"
                               value="patient"
                               checked={formData.usertype === 'patient'}
                               onChange={e => handleInputChange(e, 'usertype')} />
                        <label htmlFor="radio" className="usertype">Patient</label>

                        <input type="radio"
                               id="doctor"
                               value="doctor"
                               checked={formData.usertype === 'doctor'}
                               onChange={e => handleInputChange(e, 'usertype')} />
                        <label htmlFor="radio" className="usertype">Doctor</label>
                    </div>

                    <div className="input">
                        <label htmlFor="code" className="info">Verification Code</label>
                        <input type="text" placeholder="for doctors only" value={formData.code}  onChange={e => handleInputChange(e, 'code')}   />
                    </div>
                </div>

                <div className="submit-container">
                    <button type="button" onClick={handleFormSubmit} className="Signup">Sign UP</button>

                    <button><Link to="/Login" type="button" className="Login">Login</Link></button>
                </div>
            </form>
        </div>
    )
}
export default Signup