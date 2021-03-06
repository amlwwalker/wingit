import QtQuick 2.6
import QtQuick.Layouts 1.3
import QtQuick.Controls 2.0
import QtQuick.Controls.Material 2.0
import QtQuick.Controls.Universal 2.0
import Qt.labs.settings 1.0
import "elements"
Item {
    id: window
    width: 400
    height: 620
    visible: true
    Settings {
        id: settings
        property string style: "Default"
    }

    ColumnLayout {
    	width: parent.width
    	anchors.fill: parent
	    ToolBar {
	    	id: toolbar
	        Material.foreground: "white"
            Material.background: Material.BlueGrey
	         z: 100
            anchors.left: parent.left
            anchors.right: parent.right
            anchors.top: parent.top
	        RowLayout {
	            spacing: 20
	            anchors.fill: parent
	            ToolButton {
                    id: drawerReveal
                    visible: false
	                contentItem: Image {
	                    fillMode: Image.Pad
	                    horizontalAlignment: Image.AlignHCenter
	                    verticalAlignment: Image.AlignVCenter
                        source: titleLabel.text == "Files" ? "images/back.png" : "images/drawer.png"
    	            }
                    onClicked: {
                        if (titleLabel.text == "Files") {
                            titleLabel.text = "Contacts"
                            stackView.pop()
                        } else {
                            drawer.open()
                        }
                    }
                }

	            Label {
	                id: titleLabel
	                text: listView.currentItem ? listView.currentItem.text : "WingIt"
	                font.pixelSize: 20
	                elide: Label.ElideRight
	                horizontalAlignment: Qt.AlignHCenter
	                verticalAlignment: Qt.AlignVCenter
	                Layout.fillWidth: true
	            }
                ToolButton {
                    id: syncButton
                    contentItem: Image {
                        fillMode: Image.Pad
                        horizontalAlignment: Image.AlignHCenter
                        verticalAlignment: Image.AlignVCenter
                        source: "images/FA/teal/png/32/refresh.png"
                    }
                    ToolTip.timeout: 5000
                    ToolTip.visible: hovered
                    ToolTip.text: "Click to check for new files"
                    onClicked: {
                        QmlBridge.syncFiles()
                    }
                }
	        }
	    }
        Toast {
            //a toast that everyone can use
            id: globalToast
            x: parent.width / 10
            y: (parent.height * 4) / 5
            width: (parent.width * 4) / 5
        }

        ToolBar {
            id: footer
            Material.foreground: "white"
            Material.background: Material.BlueGrey
             z: 100
            anchors.left: parent.left
            anchors.right: parent.right
            anchors.bottom: parent.bottom
            RowLayout {
                spacing: 20
                anchors.fill: parent
                ToolButton {
                    anchors.left: parent.left
                    contentItem: Image {
                        fillMode: Image.Pad
                        horizontalAlignment: Image.AlignHCenter
                        verticalAlignment: Image.AlignVCenter
                        source: "images/FA/black/png/32/github.png"
                    }
                    onClicked: {
                        aboutDialog.open()
                    }
                }
                Label {
                    id: footerLabel
                    text: ""
                    visible: true
                    font.pixelSize: 16
                    elide: Label.ElideRight
                    horizontalAlignment: Qt.AlignHCenter
                    verticalAlignment: Qt.AlignVCenter
                    Layout.fillWidth: true
                }
                ProgressBar {
                    id: progressIndicator
                    value: 0.0
                    indeterminate: false
                    visible: false
                    z: 100
                    width: parent.width
                    anchors.horizontalCenter: parent.horizontalCenter
                    Material.accent: Material.Grey
                }

                Connections {
                    target: QmlBridge
                    //Progress bar update
                    onUpdateProcessStatus: {
                        //initialise the viewing
                        footerLabel.visible = false
                        progressIndicator.visible = true
                        progressIndicator.indeterminate = indeterminate
                        //set the progress value (only useful when determinate)
                        progressIndicator.value = c
                        if (c.toFixed(2) >=  0.98) {
                            //process complete
                            progressIndicator.visible = false
                        }
                    }
                    onSignalLogin: {
                        // receives the user who is now logged in
                        // the contacts will be loaded by the back end
                        //so all we have to do is welcome them
                        //and redirect them to the contacts page
                        drawerReveal.visible = true
                        syncButton.visible = true
                        loginButton.visible = false
                        homePageLabel.text = "Now head to your contacts to get started!"
                        homePageLabel.visible = true
                        footerLabel.text = " Welcome " + username + "!"
                        titleLabel.text = "Contacts"
                        usern.text = username
                        stackView.push("pages/_contactsPage.qml");
                    }
                    onSignalLogout: {
                        //hide menus
                        drawerReveal.visible = false
                        syncButton.visible = false
                        stackView.push("pages/_loginPage.qml");
                        usern.text = ""
                    }
                }
                ToolButton {
                    id: settingsViewer
                    anchors.right: parent.right
                    contentItem: Image {
                        fillMode: Image.Pad
                        horizontalAlignment: Image.AlignHCenter
                        verticalAlignment: Image.AlignVCenter
                        source: "qrc:/qml/images/menu.png"
                    }
                    onClicked: optionsMenu.open()

                    Menu {
                        id: optionsMenu
                        MenuItem {
                            text: "Account"
                            onTriggered: {
                                console.log("account")
                                // userPopup.open()
                                titleLabel.text = "Account/Settings"
                                stackView.push("pages/_aboutPage.qml");
                            }
                        }
                        MenuItem {
                            text: "About"
                            onTriggered: aboutDialog.open()
                        }
                        MenuItem {
                            text: "Log Out"
                            onTriggered: {
                                logoutDialog.open()
                            }
                        }
                    }
                }
            }
        }
        //content holder
	    StackView {
	        id: stackView
	        anchors.top: toolbar.bottom
	        anchors.left: parent.left
	        anchors.right: parent.right
	        anchors.bottom: footer.top
            anchors.margins: 10
              Connections {
                target: QmlBridge
                //hotloading:
                onUpdateLoader: {
                    stackView.clear()
                    stackView.push(p)
                    footerLabel.text = "detected change to: " + p
                    //loadingIndicator.visible = true
                }
              }

//animates the loader for 1 second when respawning a page for effect
	        initialItem: Pane {
	            id: pane

	            Image {
	                id: logo
	                width: pane.availableWidth / 2
	                height: pane.availableHeight / 2
	                anchors.centerIn: parent
	                anchors.verticalCenterOffset: -50
	                fillMode: Image.PreserveAspectFit
	                source: "images/wingitlogo.png"
	            }

                Button {
                    id: loginButton
                    text: "Login and get started"
                    anchors.margins: 20
                    anchors.top: logo.bottom
                    anchors.left: parent.left
                    anchors.right: parent.right
                    onClicked: function() {
                    titleLabel.text = "Login"
                        QmlBridge.createNotification("Authentication", "Log in clicked")
                        stackView.push("pages/_loginPage.qml");
                    }
                    //attempt a login when the page is loaded
                    //incase there is one in the db we can use
                    Component.onCompleted: {
                        QmlBridge.checkForUser();
                    }
                }
                Label {
                    id: homePageLabel
                    visible: false
                    anchors.margins: 20
                    anchors.top: logo.bottom
                    anchors.left: parent.left
                    anchors.right: parent.right
                    horizontalAlignment: Label.AlignHCenter
                    verticalAlignment: Label.AlignVCenter
                    wrapMode: Label.Wrap
                }

	        }
	    }
	}
    //menu
    Drawer {
        id: drawer
        width: Math.min(window.width, window.height) / 3 * 2
        height: window.height
        dragMargin: stackView.depth > 1 ? 0 : undefined

        ListView {
            id: listView
            currentIndex: -1
            anchors.fill: parent

            delegate: ItemDelegate {
                width: parent.width
                text: model.title
                highlighted: ListView.isCurrentItem

                onClicked: {
                    if (listView.currentIndex != index) {
                        listView.currentIndex = index
                        console.log("source: " + model.source)
                        stackView.push(model.source)
                        titleLabel.text = model.title
                        if (model.title == "Downloads") {
                            QmlBridge.getDownloads()
                        }
                    }
                    drawer.close()
                }
            }

            model: ListModel {
                ListElement { title: "Contacts"; source: "pages/_contactsPage.qml" }
                ListElement { title: "Search"; source: "pages/_searchPage.qml" }
                ListElement { title: "Downloads"; source: "pages/_downloadsPage.qml" }

            }

            ScrollIndicator.vertical: ScrollIndicator { }
        }
    }
    Popup {
        id: notificationPopup
        x: (window.width - width) / 2
        y: window.height / 6
        width: Math.min(window.width, window.height) / 3 * 2
        height: settingsColumn.implicitHeight + topPadding + bottomPadding
        modal: true
        focus: true
        contentItem: Text {
            id: notifcationText
            text: ""
        }
    }
    Popup {
        id: userPopup
        x: (window.width - width) / 2
        y: window.height / 6
        width: Math.min(window.width, window.height) / 3 * 2
        height: settingsColumn.implicitHeight + topPadding + bottomPadding
        modal: true
        focus: true

        contentItem: ColumnLayout {
            id: settingsColumn
            spacing: 20

            Label {
                text: "Account details"
                font.bold: true
            }

            Text { id: usern; text: QmlBridge.username; font.bold: true; }

            RowLayout {
                spacing: 10

                Button {
                    id: okButton
                    text: "Ok"
                    onClicked: {
                        userPopup.close()
                        console.log(QmlBridge.username)
                    }

                    Material.foreground: Material.primary
                    Material.background: "transparent"
                    Material.elevation: 0

                    Layout.preferredWidth: 0
                    Layout.fillWidth: true
                }
            }
        }
    }
    Popup {
        id: logoutDialog
        modal: true
        focus: true
        x: (window.width - width) / 2
        y: window.height / 6
        width: Math.min(window.width, window.height) / 3 * 2
        contentHeight: logoutColumn.height

        Column {
            id: logoutColumn
            spacing: 20

            Label {
                text: "Logout"
                font.bold: true
            }

            Label {
                width: logoutDialog.availableWidth
                text: "Logging out is really only necessary if you share a computer with others. If you want to continue, click logout below otherwise click elsewhere on the window"
                wrapMode: Label.Wrap
                font.pixelSize: 12
            }
            Button {
                id: searchButton
                text: "Logout"
                width: logoutDialog.availableWidth
                onClicked: function() {
                    console.log('logging out')
                    QmlBridge.logout()
                    logoutDialog.close()
                }
            }
        }
    }
    Popup {
        id: aboutDialog
        modal: true
        focus: true
        x: (window.width - width) / 2
        y: window.height / 6
        width: Math.min(window.width, window.height) / 3 * 2
        contentHeight: aboutColumn.height

        Column {
            id: aboutColumn
            spacing: 20

            Label {
                text: "About"
                font.bold: true
            }

            Label {
                width: aboutDialog.availableWidth
                text: "WingIt is developed in an open manner, to offer file sharing in a closed manner."
                wrapMode: Label.Wrap
                font.pixelSize: 12
            }

            Label {
                width: aboutDialog.availableWidth
                text: "All code and details are available at <a href='https://github.com/amlwwalker/got-qt'>github.com/amlwwalker/got-qt</a> or read about Wingit at <a href='https://twitter.com/amlwwalker'>twitter.com/amlwwalker</a>"
                wrapMode: Label.Wrap
                font.pixelSize: 12
                onLinkActivated: Qt.openUrlExternally(link)
            }
        }
    }
}
