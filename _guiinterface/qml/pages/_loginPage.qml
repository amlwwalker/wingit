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
            Material.foreground: Material.BlueGrey
            MouseArea {
                anchors.fill: parent
                cursorShape: Qt.PointingHandCursor
                onClicked: {
                    footerLabel.text = "Clicked on " + labelText
                    QmlBridge.addResultToContacts(labelText)
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
            text: "Use the access key you were given when you created an account, or retrieve your key at <a href='https://app-wingit.herokuapp.com/downloads/page'>Wingit</a>"
            onLinkActivated: Qt.openUrlExternally(link)
        }
        TextField {
            id: apiKeyField
            Layout.fillWidth: true
            placeholderText: "abc123"
            horizontalAlignment: Qt.AlignHCenter
            Keys.onReturnPressed: {
                footerLabel.text = "searching for " + apiKeyField.text
                QmlBridge.loginUser(apiKeyField.text)
                apiKeyField.text = ""
            }
            background: Rectangle {
                border.color: "grey"
                radius: 2
            }
        }
        Button {
            id: authButton
            text: "Login"
            Layout.fillWidth: true
            onClicked: function() {
                console.log("the user is logging in with " + apiKeyField.text)
                footerLabel.text = "logging in " + apiKeyField.text
                QmlBridge.loginUser(apiKeyField.text)
                apiKeyField.text = ""
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

                property string labelText: email
                property ListView view: listView
                property int ourIndex: index
            }
        }
    }
}
