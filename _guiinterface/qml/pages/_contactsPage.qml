import QtQuick 2.6
import QtQuick.Layouts 1.1
//import QtGraphicalEffects 1.0
import QtQuick.Controls 2.0
import QtQuick.Controls.Material 2.0
import QtQuick.Dialogs 1.0
import "../elements"
import "../images"
Pane {
    padding: 0
    property var delegateComponentMap: {
        "ItemDelegate": itemDelegateComponent
    }

    Component {
        id: itemDelegateComponent

        ItemDelegate {
            id: deleg
            //text: labelText
            width: parent.width
            Material.foreground: Material.BlueGrey
            ToolTip.timeout: 5000
            ToolTip.visible: hovered
            ToolTip.text: "Click to choose a file, or drag a file here, to send"
            MouseArea {
                anchors.fill: parent
                cursorShape: Qt.PointingHandCursor
                acceptedButtons: Qt.LeftButton | Qt.RightButton
                onClicked: {
                    console.log("file count: " + fC)
                    console.log("updating user: " + labelText)
                    QmlBridge.updatePendingUser(labelText)
                    titleLabel.text = "Files"
                    //this will work in production
                    //bit messy at the moment having both hot load and cold load
                    //working on same files
                    // stackView.push("qrc:/qml/pages/_filelistPage.qml");
                    stackView.replace("_filelistPage.qml");
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
                        // color: "#607D8B"
                    }
                    Row {
                        spacing: 5
                        Repeater {
                            model: labelText
                            Text { text: labelText; }
                        }
                    }
                }
            }

            Row {
                anchors.verticalCenter: deleg.verticalCenter
                anchors.right: deleg.right
                spacing: 10

                Text {
                    id: costText
                    anchors.verticalCenter: parent.verticalCenter
                    text: fC + " file" + (fC == 1 ? "" : "s")
                    font.pixelSize: 15
                    // color: "#607D8B"
                }
                Image {
                    id: uploadImage
                    source: mouseArea.containsMouse ? "../images/FA/teal/png/32/upload.png" : "../images/FA/bluegrey/png/32/upload.png"
                    MouseArea {
                        id: mouseArea
                        anchors.fill:parent;
                        cursorShape: Qt.PointingHandCursor
                        hoverEnabled: true
                        onClicked: {
                            fileDialog.open()
                            footerLabel.text = "Clicked on " + labelText
                        }
                    }
                }
                BusyIndicator {
                    id: loadingIndicator
                    visible: false
                    z: 100
                    width: 50
                    height: 50
                    Material.accent: Material.BlueGrey
                }
            }
            FileDialog {
                id: fileDialog
                title: "Please choose a file"
                folder: shortcuts.home
                onAccepted: {
                    console.log("You chose: " + fileDialog.fileUrls)
                    fileDialog.close()
                    //pass the current user (pending user) and the file path selected
                    //not sure how to handle multiple selections at the moment
                    //as it might be passing an array...
                    // globalToast.open()
                    // globalToast.start("uploading " + fileDialog.fileUrls.toString().replace(/^.*[\\\/]/, ''))
                    uploadImage.visible = false
                    loadingIndicator.visible = true
                    QmlBridge.updatePendingUpload(labelText, fileDialog.fileUrls)
                }
                onRejected: {
                    console.log("Canceled")
                    fileDialog.close()
                }
            }
            DropArea {
                id: drg
                anchors.fill: parent
                onDropped: {
                    uploadImage.visible = false
                    loadingIndicator.visible = true
                    var uploadResponse = QmlBridge.updatePendingUpload(labelText, drop.urls)
                    uploadImage.source = "../images/FA/bluegrey/png/32/upload.png"
                    globalToast.open()
                    globalToast.start(uploadResponse)
                }
                onEntered: {
                    //change the icon
                    uploadImage.source = "../images/FA/teal/png/32/upload.png"
                }
                onExited: {
                    //set the icon back
                    uploadImage.source = "../images/FA/bluegrey/png/32/upload.png"
                }
            }
            Connections {
                target: QmlBridge
                //indicator update
                onUpdateProcessStatus: {
                    if (c.toFixed(2) >=  0.98) {
                        //process complete
                        uploadImage.visible = true
                        loadingIndicator.visible = false
                    }
                }
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
            padding: 5
            bottomPadding: 0
            topPadding: 0
            horizontalAlignment: Qt.AlignHLeft
            text: "These are your contacts. <br><ul><li>Left click to see files from them</li><li>Click the Arrow to choose a file to send to them</li><li>You can also drag a file onto a contact to send it to them.</li></ul>"
        }
        //Button {
        //    id: button
        //    text: "Search for new contacts"
        //    anchors.margins: 10
        //    anchors.topMargin: 0
        //    anchors.top: logo.bottom
        //    anchors.left: parent.left
        //    anchors.right: parent.right
        //    onClicked: function() {
        //        globalToast.open()
        //        globalToast.start("Search for contacts")
        //    }
        //}
        ListView {
            id: listView
            Layout.fillWidth: true
            Layout.fillHeight: true
            clip: true
            model: ContactsModel
            delegate: Loader {
                id: delegateLoader
                width: listView.width
                sourceComponent: delegateComponentMap["ItemDelegate"]

                property string labelText: email
                property string fC: fileCount
                property ListView view: listView
                property int ourIndex: index

                // Can't find a way to do this in the SwipeDelegate component itself,
                // so do it here instead.
                ListView.onRemove: SequentialAnimation {
                    PropertyAction {
                        target: delegateLoader
                        property: "ListView.delayRemove"
                        value: true
                    }
                    NumberAnimation {
                        target: item
                        property: "height"
                        to: 0
                        easing.type: Easing.InOutQuad
                    }
                    PropertyAction {
                        target: delegateLoader
                        property: "ListView.delayRemove"
                        value: false
                    }
                }
            }
        }
    }
}
