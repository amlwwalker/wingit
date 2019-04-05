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
            text: labelText
            
            width: parent.width
            MouseArea {
                anchors.fill: parent
                cursorShape: Qt.PointingHandCursor
                onClicked: {
                    globalToast.open()
                    globalToast.start("opening " + labelText)
                    QmlBridge.openFile(labelText)
                }

            }
            ToolTip.timeout: 5000
            ToolTip.visible: hovered
            ToolTip.text: "Click to open the file"
        }
    }
    ColumnLayout {
        spacing: 10
        anchors.fill: parent
        anchors.topMargin: 20
        Label {
            Layout.fillWidth: true
            wrapMode: Label.Wrap
            padding: 20
            topPadding: 0
            bottomPadding: 0
            horizontalAlignment: Qt.AlignHLeft
            text: "These are files you have downloaded. <br><ul><li>Click on one to open it</li><li>Files are stored at <a href='~/.wingit/downloads/'>~/.wingit/downloads/</a></li></ul>"
        }
        Button {
            id: searchButton
            text: "Open downloads folder"
            Layout.fillWidth: true
            // horizontalAlignment: Qt.AlignHLeft
            // background: Rectangle {
            //     // color: "#EEEEEE"
            //     radius: 2
            //     // border.color: "grey"
            // }
            // Material.background: "#BEEEFF"
            onClicked: function() {
                globalToast.open()
                globalToast.start("opening downloads folder")
                QmlBridge.openDowloadsDirectory()
            }
        }
        ListView {
            id: listView
            Layout.fillWidth: true
            Layout.fillHeight: true
            clip: true
            spacing: 2
            model: DownloadsModel

            section.property: "type"

            delegate: Loader {
                id: delegateLoader
                width: listView.width
                sourceComponent: delegateComponentMap["ItemDelegate"]

                property string labelText: filePath
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
