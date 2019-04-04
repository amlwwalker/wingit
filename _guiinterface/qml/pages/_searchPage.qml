import QtQuick 2.6
import QtQuick.Layouts 1.1
import QtQuick.Controls 2.0
import QtQuick.Controls.Material 2.0
Pane {
    padding: 0
    property var delegateComponentMap: {
        "ItemDelegate": itemDelegateComponent
    }

    Component {
        id: itemDelegateComponent
        ItemDelegate {
            id: deleg
            // property int indexOfThisDelegate: index
            // text: labelText
            width: parent.width
            Material.foreground: Material.BlueGrey
            MouseArea {
                anchors.fill: parent
                cursorShape: Qt.PointingHandCursor
                onClicked: {
                    footerLabel.text = "adding " + labelText
                    console.log("index added (and removed from this list)", ourIndex)
                    globalToast.open()
                    globalToast.start("New contact added!")
                    QmlBridge.addResultToContacts(labelText, ourIndex)
                }

            }


            Row {
                anchors.verticalCenter: deleg.verticalCenter
                spacing: 10

                Column {
                    anchors.verticalCenter: parent.verticalCenter

                    Text {
                        text: labelText
                        font.pixelSize: 15
                        color: "#607D8B"
                    }
                    Row {
                        spacing: 5
                        Repeater {
                            model: labelText
                            Text { text: labelText; color: "#607D8B" }
                        }
                    }
                }
            }
            Row {
                anchors.verticalCenter: deleg.verticalCenter
                anchors.right: deleg.right
                spacing: 10
                Image {
                    id: addContactImage
                    source: "../images/FA/bluegrey/png/32/plus-circle.png"
                    // source: mouseArea.containsMouse ? "../images/FA/teal/png/32/plus-circle.png" : "../images/FA/bluegrey/png/32/plus-circle.png"
                    // MouseArea {
                    //     id: mouseArea
                    //     anchors.fill:parent;
                    //     cursorShape: Qt.PointingHandCursor
                    //     hoverEnabled: true
                    //     onClicked: {
                    //         // SearchModel.remove(indexOfThisDelegate)
                    //         // fileDialog.open()
                    //         // footerLabel.text = "Clicked on " + labelText
                    //     }
                    // }
                }
            }


            ToolTip.timeout: 5000
            ToolTip.visible: hovered
            ToolTip.text: "Click to add contact to your contacts"
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
            text: "Search for new contacts here. Click on a search result to add them to your contacts"
        }
        TextField {
            id: searchContactField
            Layout.fillWidth: true
            placeholderText: "James Bond"
            horizontalAlignment: Qt.AlignHCenter
            Keys.onReturnPressed: {
                console.log("the user is searching for " + searchContactField.text)
                footerLabel.text = "searching for " + searchContactField.text
                searchContactField.text = ""
            }
            background: Rectangle {
                border.color: "grey"
                radius: 2
            }
        }
        Button {
            id: searchButton
            text: "Search"
            Layout.fillWidth: true
            Material.background: "#BEEEFF"
            onClicked: function() {
                footerLabel.text = "searching for " + searchContactField.text
                QmlBridge.searchFor(searchContactField.text)
                searchContactField.text = ""
            }
        }
        ListView {
            id: listView
            Layout.fillWidth: true
            Layout.fillHeight: true
            clip: true
            spacing: 2
            model: SearchModel
            delegate: Loader {
                id: delegateLoader
                width: listView.width
                sourceComponent: delegateComponentMap["ItemDelegate"]

                //this is getting the new contact name from 'email'
                //this is a field on a person
                //we can use this principle elsewhere
                property string labelText: email
                property ListView view: listView
                property int ourIndex: index
            }
        }
    }
}
