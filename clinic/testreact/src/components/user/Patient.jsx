
import './Patient.css';

import {Link, useLocation} from "react-router-dom";
import axios from "axios";
import React, { useState, useEffect } from 'react';

function Patient() {
    const [showCreateAppointment, setShowCreateAppointment] = useState(false);
    const [showUpdateAppointment, setShowUpdateAppointment] = useState(false);
    const [showCancelAppointment, setShowCancelAppointment] = useState(false);
    const [UpdateDoctor, setUpdateDoctor] = useState(false);
    const [UpdateSlot, setUpdateSlot] = useState(false);

    const handleCreateClick = () => {
        setShowCreateAppointment(true);
        setShowUpdateAppointment(false);
        setShowCancelAppointment(false);
    }

    const handleUpdateClick = () => {
        setShowCreateAppointment(false);
        setShowUpdateAppointment(true);
        setShowCancelAppointment(false);
    }

    const handleCancelClick = () => {
        setShowCreateAppointment(false);
        setShowUpdateAppointment(false);
        setShowCancelAppointment(true);
    }

    const username = localStorage.getItem('username');
    const [formData, setFormData] = useState({ patientName:'',doctorName:'', date: '',time:''}); // Set the default user type to 'patient'
    const [error, setError] = useState('');
    const [formDataaaa, setFormDataaa] = useState({ patientName:'',doctorName:'', date: '',time:''}); // Set the default user type to 'patient'

    const handleFormSubmit = (event) => {
        event.preventDefault();  //Prevent page reload on form submission
        let dataToSend = {
            doctorName: formData.doctorName,
            date: formData.date,
            time: formData.time,
        };

        axios.post(`http://localhost:12345/patient/CancelReservation`, dataToSend)
            .then(response => {
                console.log(response.data);
                setError(" appointment deleted Successfully");
            })
            .catch(error => {
                console.error(error);
                setError("Please Try Again");
            });
    }

    /////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    const [appointments, setAppointments] = useState([]);
    //show the table of all patient's appointments
    useEffect(() => {
        async function fetchAppointments() {
            try {
                const response = await fetch(`http://localhost:12345/patient/AllReservation`); // Replace with your API endpoint
                if (!response.ok) {
                    throw new Error('Network response was not ok');
                }
                const data = await response.json();
                const filteredAppointments = data.filter(appointment => appointment.patientName === username);
                setAppointments(filteredAppointments);
            } catch (error) {
                // Handle error, e.g., display an error message
                console.error('Error fetching appointments:', error);
            }
        }
        fetchAppointments();
    }, []);
    ////////////////////////////////////////////////////////////////////////////////////////////////////////
    const [availableSlots, setAvailableSlots] = useState([]);
    //get all slots to reserve appointment
    const handleShowSlots = async () => {
        try {
            const response = await fetch(`http://localhost:12345/patient/Getslot`, {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',
                },
            });
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            const data = await response.json();
            setAvailableSlots(data);

        } catch (error) {
            console.error('Error fetching available slots:', error);
            // Handle error, e.g., display an error message
        }
    };


    /////////////////////////////////////////////////////////////////////////////////////
    //reserve appointment
    const handleFormCreate = (event) => {
        event.preventDefault();  //Prevent page reload on form submission
        let dataToSend = {
            patientName:username,
            doctorName: formDataaaa.doctorName,
            date: formDataaaa.date,
            time: formDataaaa.time,
        };

        axios.post(`http://localhost:12345/patient/ReserveAppointment`, dataToSend)
            .then(response => {
                console.log(response.data);
                setError(" appointment created Successfully");
            })
            .catch(error => {
                console.error(error);
                setError("Please Try Again");

            });
    }

    /////////////////////////////////////////////////////////////////////////////////////
    //update doctor
    const [formDataDoc, setFormDataDoc] = useState({name:'',oldDate:'',oldTime:''}); // Set the default user type to 'patient'

    const [errorDoc, setErrorDoc] = useState('');

    const [formDataDocNew, setFormDataDocNew] = useState({newName:'', date: '',time:''}); // Set the default user type to 'patient'

    const handleFormSubmitDoc = (event) => {
        event.preventDefault();  // Prevent page reload on form submission
        let dataToSend = {
            name: formDataDoc.name,
            oldDate:formDataDoc.oldDate,
            oldTime:formDataDoc.oldTime,
            newName:formDataDocNew.newName,
            date: formDataDocNew.date,
            time: formDataDocNew.time,
        };
        axios.post(`http://localhost:12345/patient/UpdateReservation/Doctor`, dataToSend)
            .then(response => {
                console.log(response.data);
                setErrorDoc(" doctor changed Successfully");
            })
            .catch(errorDoc => {
                console.error(errorDoc);
                setErrorDoc('Please Try Again');

            });
    }
//////////////////////////////////////////////////////////////////////////////////////////////////////

    const [formDataSlot, setFormDataSlot] = useState({doctorName:'',oldDate:'',oldTime:'',date: '',time:''}); // Set the default user type to 'patient'

    const [errorSlot, setErrorSlot] = useState('');

    const handleFormSubmitSlot = (event) => {
        event.preventDefault();  // Prevent page reload on form submission
        let dataToSend = {
            doctorName: formDataSlot.doctorName,
            oldDate:formDataSlot.oldDate,
            oldTime:formDataSlot.oldTime,
            date: formDataSlot.date,
            time: formDataSlot.time,
        };
        axios.post(`http://localhost:12345/patient/UpdateReservation/Slot`, dataToSend)
            .then(response => {
                console.log(response.data);
                setErrorSlot('Slot changed Successfully');
            })
            .catch(errorDoc => {
                console.error(errorDoc);
                setErrorSlot('Please Try Again');
            });
    }

    //////////////////////////////////////////////////////////////////////////////////////////////////


    return (
        <div className="container">
            <div className="header">
                <div className="text">Hello Patient {username} </div>
                <div className="underline"> </div>
            </div>
            <div>
                <h2>Your Appointments</h2>
                <table>
                    <thead>
                    <tr>
                        <th>Patient</th>
                        <th>Doctor</th>
                        <th>Date</th>
                        <th>Time</th>
                    </tr>
                    </thead>
                    <tbody>
                    {appointments.map((appointment) => (
                        <tr key={appointment.ID}>
                            <td>{appointment.patientName}</td>
                            <td>{appointment.doctorName}</td>
                            <td>{appointment.date}</td>
                            <td>{appointment.time}</td>
                        </tr>
                    ))}
                    </tbody>

                </table>
            </div>


            <button className="button" type='submit' onClick={handleShowSlots}>show available slots</button>
            <div>
                <h2>Display the available slots</h2>
                <ul>
                    {availableSlots.map((slot) => (
                        <li key={slot.id}>
                            {/* Display the available slots information here */}
                            Doctor: {slot.name}            ,
                            Date: {slot.date}              ,
                            Time: {slot.time}
                            <p>-----------------------</p>
                        </li>
                    ))}
                </ul>
            </div>

            <div className="cancelApp">
                <button className="button" onClick={handleCreateClick}>Create Appointment</button>
                <button className="button" onClick={handleUpdateClick}>Update Appointment</button>
                <button className="button" onClick={handleCancelClick}>Cancel Appointment</button>
            </div>




            {/*/////////////////////////////////////////////////////////////////////////////////////////*/}
            {showCreateAppointment && (
                <div>
                    <div className="cancelApp">
                        <input
                            type="text"
                            name="doctorName1"
                            placeholder="enter doctor"
                            value={formDataaaa.doctorName} onChange={e => setFormDataaa({ ...formDataaaa, doctorName: e.target.value })}
                        />
                        <input
                            type="text"
                            name="date1"
                            placeholder="enter date"
                            value={formDataaaa.date} onChange={e => setFormDataaa({ ...formDataaaa, date: e.target.value })}
                        />
                        <input
                            type="text"
                            name="time1"
                            placeholder="enter time"
                            value={formDataaaa.time} onChange={e => setFormDataaa({ ...formDataaaa, time: e.target.value })}
                        />

                    </div>
                    <button className="button2" type='submit' onClick={handleFormCreate}>Reserve</button>
                </div>
            )}

            {/*///////////////////////////////////////////////////////////////////////////////////////////////*/}
            {showUpdateAppointment && (
                <div>
                    ---------------------------------------------------------------------------------------------
                    <div className="cancelApp">
                        <button className="button" onClick={() => setUpdateDoctor(true)}>Update Doctor</button>
                        <button className="button" onClick={() => setUpdateSlot(true)}>Update Slot</button>
                    </div>
                    {UpdateDoctor && (
                        <div>
                            ---------------------------------------------------------------------------------------------

                            <h2>Update Doctor</h2>
                            <div className="cancelApp">
                                <input
                                    type="text"
                                    name="OdoctorName"
                                    placeholder="enter current doctor"
                                    value={formDataDoc.name} onChange={e => setFormDataDoc({ ...formDataDoc, name: e.target.value })}
                                />
                                <input
                                    type="text"
                                    name="Odate"
                                    placeholder="enter current date"
                                    value={formDataDoc.oldDate} onChange={e => setFormDataDoc({ ...formDataDoc,oldDate: e.target.value })}
                                />
                                <input
                                    type="text"
                                    name="Otime"
                                    placeholder="enter current time"
                                    value={formDataDoc.oldTime} onChange={e => setFormDataDoc({ ...formDataDoc,oldTime: e.target.value })}
                                />
                            </div>

                            <div className="cancelApp">
                                <input
                                    type="text"
                                    name="NdoctorName"
                                    placeholder="enter new doctor"
                                    value={formDataDocNew.newName} onChange={e => setFormDataDocNew({ ...formDataDocNew,newName: e.target.value })}
                                />

                                <input
                                    type="text"
                                    name="Ndate"
                                    placeholder="enter new date"
                                    value={formDataDocNew.date} onChange={e => setFormDataDocNew({ ...formDataDocNew, date: e.target.value })}
                                />
                                <input
                                    type="text"
                                    name="Ntime"
                                    placeholder="enter new time"
                                    value={formDataDocNew.time} onChange={e => setFormDataDocNew({ ...formDataDocNew, time: e.target.value })}
                                />
                            </div>
                            <button className="button2" type='submit' onClick={handleFormSubmitDoc}>{'\n'} Update Doctor</button>
                        </div>
                    )}
                    {UpdateSlot && (
                        <div>
                            ---------------------------------------------------------------------------------------------

                            <h2>Update Slot</h2>
                            <div className="cancelApp">
                                <input
                                    type="text"
                                    name="doctorName"
                                    placeholder="enter doctor"
                                    value={formDataSlot.doctorName} onChange={e => setFormDataSlot({ ...formDataSlot, doctorName: e.target.value })}
                                />
                                <input
                                    type="text"
                                    name="date"
                                    placeholder="enter date"
                                    value={formDataSlot.oldDate} onChange={e => setFormDataSlot({ ...formDataSlot, oldDate: e.target.value })}
                                />
                                <input
                                    type="text"
                                    name="time"
                                    placeholder="enter time"
                                    value={formDataSlot.oldTime} onChange={e => setFormDataSlot({ ...formDataSlot, oldTime: e.target.value })}
                                />
                            </div>
                            <div className="cancelApp">

                                <input
                                    type="text"
                                    name="Ndate"
                                    placeholder="enter new date"
                                    value={formDataSlot.date} onChange={e => setFormDataSlot({ ...formDataSlot, date: e.target.value })}
                                />
                                <input
                                    type="text"
                                    name="Ntime"
                                    placeholder="enter new time"
                                    value={formDataSlot.time} onChange={e => setFormDataSlot({ ...formDataSlot, time: e.target.value })}
                                />
                            </div>
                            <button className="button2" type='submit' onClick={handleFormSubmitSlot}>{'\n'} Update Slot</button>
                        </div>

                    )}

                </div>
            )}
            {errorDoc && <div className="error-message header"  style={{ color: 'green',fontSize:'30px' }}>{errorDoc}</div>}
            {errorSlot && <div className="error-message header"  style={{ color: 'green',fontSize:'30px' }}>{errorSlot}</div>}


            {/*///////////////////////////////////////////////////////////////////////////////////////////////*/}
            {error && <div className="error-message header "  style={{ color: 'green',fontSize:'30px' }}>{error}</div>}
            {showCancelAppointment && (
                <div>
                    ---------------------------------------------------------------------------------------------

                    <h2>Cancel Appointment</h2>
                    <div class="cancelApp">
                        <input
                            type="text"
                            name="doctorName"
                            placeholder="enter doctor"
                            value={formData.doctorName} onChange={e => setFormData({ ...formData, doctorName: e.target.value })}
                        />
                        <input
                            type="text"
                            name="date"
                            placeholder="enter date"
                            value={formData.date} onChange={e => setFormData({ ...formData, date: e.target.value })}
                        />
                        <input
                            type="text"
                            name="time"
                            placeholder="enter time"
                            value={formData.time} onChange={e => setFormData({ ...formData, time: e.target.value })}
                        />
                    </div>
                    <button class="button2" type='submit' onClick={handleFormSubmit}>{'\n'} Cancel Appointment</button>
                </div>
            )}

        </div>
    )
}
export default Patient