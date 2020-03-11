package main

import(
    "./config"
    "./elevio"
    . "./elevatortypes"
    //"./utils"
    "./elevator"
    "./network/bcast"
    //"./network/broadcaster"
    //"./network/localip"
    "./fsm"

    "fmt"
    "strconv"
)

func main(){
    LOCAL_ID := "8080"
    globalState := elevator.GlobalElevatorInit(config.N_FLOORS, config.N_BUTTONS, LOCAL_ID)
    fmt.Println(globalState.ID)

    //hardware
    localAddress := "localhost:" + strconv.Itoa(config.SERVER_PORT)
	elevio.Init(localAddress, config.N_FLOORS)
	buttonOrderC := make(chan ButtonEvent)
	floorEventC := make(chan int)
	go elevio.PollButtons(buttonOrderC)
	go elevio.PollFloorSensor(floorEventC)
    fmt.Println(localAddress)


    doorTimerFinishedC := make(chan bool)
	doorTimerStartC := make(chan bool, 10) //Buffered to avoid blocking
	go fsm.InitDoorTimer(doorTimerFinishedC, doorTimerStartC, config.DOOR_OPEN_TIME)


    // Local fsm
    localStateUpdateC := make(chan Elevator)
    updateFSMRequestsC := make(chan [][]bool, 10)
    fsm.InitFsm(
		localStateUpdateC,
		updateFSMRequestsC,
		floorEventC,
		doorTimerFinishedC,
		doorTimerStartC,
		LOCAL_ID,
		config.N_FLOORS,
		config.N_BUTTONS)

    networkTxC := make(chan GlobalElevator)
    networkRxC := make(chan GlobalElevator)
    go bcast.Transmitter(config.BROADCAST_PORT, networkTxC)
    go bcast.Receiver(config.BROADCAST_PORT, networkRxC)

}
