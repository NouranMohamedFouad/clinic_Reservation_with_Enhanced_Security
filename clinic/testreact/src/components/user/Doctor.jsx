import React, {useState, Fragment, useEffect} from "react";
// import {nanoid} from 'nanoid';
import './Doctor.css';
// import data from "../../mock-data.json"

import axios from "axios";
import { useLocation } from 'react-router-dom';


const Doctor=() =>{




    const searchParams = new URLSearchParams(useLocation().search);
    const username = localStorage.getItem('username');
    const [formData, setFormData] = useState({ name:'', date: '',time:'', isAvailable:true}); // Set the default user type to 'patient'
    const [error, setError] = useState('');


    const handleFormSubmit = (event) => {
        event.preventDefault();  // Prevent page reload on form submission
        let dataToSend = {
            name: username,
            date: formData.date,
            time: formData.time,
            isAvailable: formData.isAvailable
        };
        const token = localStorage.getItem('token');
        // Check if token exists
        if (token) {
            // Add token to request headers
            const headers = {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json'
            };


            axios.post(`http://localhost:12345/doctor/SetSchudule`, dataToSend , {headers})
                .then(response => {
                    console.log(response.data);
                    setError(" Slot Added Successfully");
                })
                .catch(error => {
                    console.error(error);
                 if (error.response.status === 401) {
                        setError('unauthorized ,please go login again!HEHE ');
                    }
                 else {
                     setError("Please Try Again");
                 }
                });
        }
    }

    ////////////////////////////////////////////////////////////////////////////////////////////
    const [slotss, setSlots] = useState([]);
    //show the table of all patient's appointments
    useEffect(() => {
        async function fetchAppointments() {
            try {
                const response = await fetch(`http://localhost:12345/doctor/AllSlots`); // Replace with your API endpoint
                if (!response.ok) {
                    throw new Error('Network response was not ok');
                }
                const data = await response.json();
                const filteredSlots = data.filter(slot => slot.name === username);
                setSlots(filteredSlots);
            } catch (error) {
                // Handle error, e.g., display an error message
                console.error('Error fetching appointments:', error);
            }
        }
        fetchAppointments();
    }, []);
    ////////////////////////////////////////////////////////////////////////////////////////////////////////




    return (
        <div className="container">
            <div className="header">
                <div className="text">Hello Doctor {username} </div>
                <div className="underline"></div>
            </div>



            {/*<div className="slotat">*/}
            {/*<h2>My Slots</h2>*/}
            {/*</div>*/}
            <div>
                <h2>Your Slots</h2>
                <table>
                    <thead>
                    <tr>
                        <th>Doctor</th>
                        <th>Date</th>
                        <th>Time</th>
                    </tr>
                    </thead>
                    <tbody>
                    {slotss.map((slot) => (
                        <tr key={slot.ID}>
                            <td>{slot.name}</td>
                            <td>{slot.date}</td>
                            <td>{slot.time}</td>
                        </tr>
                    ))}
                    </tbody>

                </table>
            </div>

            <div className="slotat">


                <h2>Create New Slot</h2>
                <form >

                    <input
                        type="text"
                        name="date"
                        placeholder="choose date"
                        value={formData.date} onChange={e => setFormData({ ...formData, date: e.target.value })}
                    />

                    <input type="text"
                           name="time"
                           placeholder="choose time"
                           value={formData.time} onChange={e => setFormData({ ...formData, time: e.target.value })}
                    />


                    <button type='submit' onClick={handleFormSubmit}>Add Slot</button>

                </form>
            </div>
            {error && <div className="error-message header "  style={{ color: 'green',fontSize:'35px' }}>{error}</div>}
        </div>

    )
}
export default Doctor