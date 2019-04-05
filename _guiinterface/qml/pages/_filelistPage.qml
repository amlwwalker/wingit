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
            
            // width: 100
            
            Text {
                text: prunedText
                // color: "#607D8B"
                anchors.verticalCenter: parent.verticalCenter
            }
            MouseArea {
                anchors.fill: parent
                cursorShape: Qt.PointingHandCursor
                onClicked: {
                    globalToast.open()
                    globalToast.start("Downloading " + labelText + "...")
                    QmlBridge.beginFileDownload(labelText)
                }

            }
            ToolTip.timeout: 5000
            ToolTip.visible: hovered
            ToolTip.text: "Click to download the file"
            Rectangle{
                id: infoArea
                color: "transparent"
                anchors.verticalCenter: parent.verticalCenter
                anchors.right: parent.right
                anchors.rightMargin: 45
                Text{
                    id: infoText
                    text: (labelFileSize / 1024 / 1024).toFixed(2) + " MB"
                    anchors.centerIn: parent
                    // color: "#607D8B"
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
            padding: 20
            topPadding: 0
            horizontalAlignment: Qt.AlignHLeft
            text: "These are files waiting for you. <br><ul><li>Click on one to download it</li><li>A progress bar will appear at the bottom</li><li>Once it has finished, you can find a list of your downloaded files in the <b>Downloads</b> tab</li><li>Note, they will still appear here while they are available for download</li></ul>"
        }
        ListView {
            id: listView
            Layout.fillWidth: true
            Layout.fillHeight: true
            clip: true
            spacing: 2
            model: FilesModel
            //ListModel {
            //    ListElement { type: "ItemDelegate"; text: "domination.pdf" }
            //    ListElement { type: "ItemDelegate"; text: "megolomania.exe" }
            //    ListElement { type: "ItemDelegate"; text: "bankvault.md" }
            //}

            section.property: "type"

            delegate: Loader {
                id: delegateLoader
                width: listView.width
                sourceComponent: delegateComponentMap["ItemDelegate"]

                property string labelText: filePath
                property string prunedText: filePath.substr(0,19-1)+(filePath.length>19?' ... ':'')+(filePath.length>25?filePath.substr(filePath.length-5,filePath.length-1):'');
                property string labelFileSize: fileSize
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
