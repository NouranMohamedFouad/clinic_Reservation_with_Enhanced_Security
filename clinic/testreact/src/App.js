

import { BrowserRouter, Route, Routes } from 'react-router-dom';
import './App.css';
import Signup from './components/login/Signup.jsx';
import Login from './components/login/Login.jsx';
import Patient from './components/user/Patient.jsx';
import Doctor from './components/user/Doctor.jsx';

function App() {
    return (
        <div>
            <BrowserRouter>
                <Routes>
                    <Route index element={<Login/>}/>
                    <Route path='/Login' element={<Login/>}/>
                    <Route path='/Signup' element={<Signup/>}/>
                    <Route path='/Doctor' element={<Doctor/>}/>
                    <Route path='/Patient' element={<Patient/>}/>
                </Routes>
            </BrowserRouter>
        </div>

    );
}

export default App;
