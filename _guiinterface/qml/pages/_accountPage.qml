import QtQuick 2.6
import QtQuick.Layouts 1.1
import QtQuick.Controls 2.0
import QtQuick.Controls.Material 2.0
Pane {
    padding: 0
    // property var delegateComponentMap: {
    //     "ItemDelegate": itemDelegateComponent
    // }

    ColumnLayout {
        spacing: 10
        anchors.fill: parent
        anchors.topMargin: 10
        Row {
            Column {
                Label {
                    Layout.fillWidth: true
                    wrapMode: Label.Wrap
                    padding: 10
                    topPadding: 0
                    bottomPadding: 0
                    horizontalAlignment: Qt.AlignHLeft
                    text: "Account and Settings"
                }
            }
        }
        ColumnLayout {
            spacing: 10
            anchors.fill: parent
            anchors.topMargin: 20

            Label {
                Layout.fillWidth: true
                wrapMode: Label.Wrap
                horizontalAlignment: Qt.AlignHCenter
                padding: 20
                topPadding: 0
                text: "View your account details, change settings and view your file and contact data directly (advanced)"
            }
            Label {
                Layout.fillWidth: true
                wrapMode: Label.Wrap
                // horizontalAlignment: Qt.AlignHCenter
                topPadding: 10
                font.pixelSize: 18
                // topPadding: 0
                text: "User details"
            }
            Column{
                spacing: 20
                Row {
                    Text {
                        text: "Name:"
                        font.bold: true
                        rightPadding: 10
                    }
                    TextEdit {
                        text: QmlUser.name
                        selectByMouse: true
                        readOnly: true
                    }
                }
                Row {
                    Text {
                        text: "Email: "
                        font.bold: true
                        rightPadding: 10
                    }
                    TextEdit {
                        text:QmlUser.email
                        selectByMouse: true
                        readOnly: true
                    }
                }
                Row {
                    Text {
                        text: "API Key: "
                        font.bold: true
                        rightPadding: 10
                    }
                    TextEdit {
                        text: QmlUser.apiKey
                        selectByMouse: true
                        readOnly: true
                    }
                }
            }
        }
        ColumnLayout {
            Label {
                Layout.fillWidth: true
                wrapMode: Label.Wrap
                // horizontalAlignment: Qt.AlignHCenter
                topPadding: 20
                font.pixelSize: 18
                // topPadding: 0
                text: "Settings"
            }
            Button {
                id: disableSyncButton
                text: QmlBridge.toggleState() ? "Auto Sync On (disable)" : "Auto Sync Off (enable)"
                Layout.fillWidth: true
                property bool isClicked: false
                // background: Rectangle {
                //     color: "#EEEEEE"
                //     radius: 2
                //     border.color: "grey"
                // }
                // Material.background: "#BEEEFF"
                onClicked: function() {
                    // isClicked = !isClicked
                    if (QmlBridge.toggleAutoSync()) {
                        disableSyncButton.text = "Auto Sync Off (enable)"
                        // console.log("sync state: ", QmlBridge.toggleAutoSync())
                        // root.syncButtonIcon.source = "../images/FA/black/png/32/refresh.png"
                    } else {
                        disableSyncButton.text = "Auto Sync on (Disable)"
                        // console.log("sync state: ", QmlBridge.toggleAutoSync())
                    }
                }
                onPressedChanged: {
                    if ( pressed ) {
                        disableSyncButton.background.color = "#DDDDDD"
                    } else {
                        disableSyncButton.background.color = "#EEEEEE"
                    }
                }
            }
        }
        ColumnLayout {
            Label {
                Layout.fillWidth: true
                wrapMode: Label.Wrap
                // horizontalAlignment: Qt.AlignHCenter
                topPadding: 20
                font.pixelSize: 18
                // topPadding: 0
                text: "Api Endpoints"
            }
            ListView {
                id: apiListView
                Layout.fillWidth: true
                Layout.fillHeight: true
                clip: true
                spacing: 2
                model: ListModel {
                    ListElement { title: "Retrieve Files"; url: "https://app-wingit.herokuapp.com/files/list?id_token=" }
                    ListElement { title: "Retrieve Keys"; url: "https://app-wingit.herokuapp.com/keys/retrieve?id_token=" }
                }

                section.property: "type"

                delegate: ItemDelegate {
                    width: parent.width
                    text: model.title

                    onClicked: {
                        console.log("opening " + model.url)
                        Qt.openUrlExternally(model.url + QmlUser.apiKey)
                    }
                }
            }
        }
    }
}
