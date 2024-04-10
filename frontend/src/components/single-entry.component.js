import React, {useState, useEffect} from "react";
import 'bootstrap/dist/css/bootstrap.css'
import {Button, Card, Row, Col} from 'react-bootstrap'
const Entry = ({entryData, setChangeIngredient, deleteSingleEntry, setChangeEntry}) => {
    return (
        <Card>
            <Row>
                <Col>Dish:{entryData !== undefined && entryData.dish}</Col>
                <Col>Ingredients:{entryData !== undefined && entryData.ingredients}</Col>
                <Col>Calories:{entryData !== undefined && entryData.calories}</Col>
                <Col>Fat:{entryData !== undefined && entryData.fat}</Col>
                <Col><Button onClick={() => deleteSingleEntry(entryData._id)}>Delete Entry</Button></Col>
                <Col><Button onClick={() => ChangeIngredient()}>Change Ingredients</Button></Col>
                <Col><Button onClick={() => ChangeEntry()}>Change Entry</Button></Col>
            </Row>
        </Card>
    )
    function ChangeIngredient() {
        setChangeIngredient(
            {
                "change": true,
                "id":   entryData._id
            }
        )
    }
    function ChangeEntry() {
        setChangeEntry(
            {
                "change":   true,
                "id"    :   entryData._id
            }
        )
    }
}

export default Entry