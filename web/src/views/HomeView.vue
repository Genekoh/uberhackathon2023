<template>
  <header>
    <h2 class="sub-header"><b>Lets Go Places</b></h2>

    <div class="car-select">
      <div class="column">
        <h3>FROM:</h3>
        <input type="text" placeholder="Enter Latitude" name="q" />
        <input type="text" placeholder="Enter Longitude" name="q" />
      </div>
      <div class="column">
        <h3>TO:</h3>
        <input type="text" placeholder="Enter Latitude" name="q" />
        <input type="text" placeholder="Enter Longitude" name="q" />
      </div>
      <div class="column">
        <h3>Type:</h3>
        <select class="type-select">
          <option value="0">Any Type:</option>
          <option value="1">SUV</option>
          <option value="2">Purpose-Built Wheelchair cars</option>
          <option value="3">Carpool</option>
          <option value="4">Bike</option>
          <option value="5">luxury Cars</option>
          <option value="6">Vans</option>
          <option value="7">Electrics</option>
          <option value="8">People Carriers</option>
        </select>
      </div>

      <div class="column">
        <h3>Budget:</h3>
        <select class="budget-select">
          <option value="0">Any Budget:</option>
          <option value="1">Up to ฿80</option>
          <option value="2">Up to ฿100</option>
          <option value="3">Up to ฿200</option>
          <option value="4">Up to ฿300</option>
          <option value="5">Up to ฿400</option>
          <option value="6">Up to ฿500</option>
          <option value="7">Up to ฿600</option>
          <option value="8">Up to ฿700</option>
          <option value="9">Up to ฿800</option>
          <option value="10">>฿1000</option>
        </select>
      </div>

      <button type="submit" class="GO-submit">GO</button>
    </div>
  </header>
  <div class="stuff">
    <div class="body1-header">
      <h1>Popular <br />Features</h1>
    </div>
    <div class="body1-icons body-content">
      <i class="fa-solid fa-car-on" style="color: #000000"></i>
      <h2 class="body1-txt center-text">Carpool</h2>
      <p class="body1-txt2 center-text">Split your rides with other people</p>
    </div>
    <div class="body1-icons body-content">
      <i
        class="fa-solid fa-handshake-angle center-text"
        style="color: #000000"
      ></i>
      <h2 class="body1-txt center-text">Assistance Riding</h2>
      <p class="body1-txt2 center-text">
        Call drivers with expertise around medical fields
      </p>
    </div>
    <div class="body1-icons body-content">
      <i class="fa-solid fa-map-location-dot" style="color: #000000"></i>
      <h2 class="body1-txt center-text">Public Navigator</h2>
      <p class="body1-txt2 center-text">
        Find public transports around your area
      </p>
    </div>
  </div>
</template>

<script></script>

<style scoped>
header {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: space-between;
  height: 16rem;
  padding: 0rem 2rem 3rem 2rem;
  background-color: black;
  border-bottom-left-radius: 2rem;
  border-bottom-right-radius: 2rem;
}

.sub-header {
  width: 80%;
  display: flex;
  justify-content: flex-start;
  color: white;
  font-size: 2rem;
  padding-top: 1rem;
}

.car-select {
  display: flex;
  align-items: center;
  justify-content: space-around;
  width: 80%;
  height: 8rem;
  border-radius: 20px;
  background-color: rgba(245, 245, 245, 1);
}

.column {
  height: 60%;
  display: flex;
  flex-direction: column;
}

.lat-long2 {
  display: flex;
  flex-direction: column;
}

.type-select {
  border-radius: 10px;
  border-color: grey;
}

.budget-select {
  border-radius: 10px;
  border-color: grey;
}

.GO-submit {
  width: 10%;
  height: 60%;
  border: 0;
  outline: none;
  border-radius: 10px;
  background-color: black;
  color: white;
  font-size: 24px;
  cursor: pointer;
  transition: box-shadow 0.2s;
}

.GO-submit:hover {
  opacity: 0.8;
}

.GO-submit:active {
  opacity: 0.5;
}

.stuff {
  padding-top: 3.5rem;
  height: 100%;
  display: flex;
  align-content: center;
  flex-wrap: wrap;
  justify-content: space-around;
}

.body-content {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.body1-header h1 {
  font-size: 4rem;
}

.body1-txt {
  color: black;
}

.body1-txt2 {
  color: black;
  font-size: 0.8rem;
}

.body1-icons {
  width: 10%;
  height: 50%;
  display: flex;
  flex-direction: column;
}

bod1-buttons {
  width: 10%;
  height: 50%;
}

.center-text {
  text-align: center;
}

.fa-solid {
  font-size: 6rem;
}
</style>
